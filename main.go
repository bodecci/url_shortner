package main

import (
	"fmt"
	"html"
	"math/big"
	"net/http"
	"sync"
)

var (
	shortURLs     = make(map[string]string) // Map to store short URLs and their corresponding long URLs
	urlCounter    = new(big.Int)            // Coubnter to generate unique URL IDs
	urlLock       = sync.RWMutex{}          // // Mutex to synchronize access to the shortURLs map
	base62Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// redirects from shortURL to the original URL
func redirectToLongURL(w http.ResponseWriter, r *http.Request) {
	short := r.URL.Path[1:] //  extracts the short URL from the request's path, remove the leading slash from the URL

	//check if the short URL exists in the map
	urlLock.RLock()
	longURL, exists := shortURLs[short]
	urlLock.Unlock()

	if exists {
		http.Redirect(w, r, longURL, http.StatusMovedPermanently) // redirect to long URL if found
	}
	if !exists {
		http.NotFound(w, r) // the short URL not found in the map
	}
}

// shortens long URL passed in
func shortenURL(w http.ResponseWriter, r *http.Request) {
	longURL := r.FormValue("url") // this retrieves the long URL from the request form data with the key url

	if longURL == "" {
		http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
		return
	}

	// generate a unique URL ID and encode it in base62
	shortURL := base62Encode(urlCounter)
	shortURLs[shortURL] = longURL
	urlCounter.Add(urlCounter, big.NewInt(1))
	urlLock.Unlock()

	shortURL = fmt.Sprintf("http://localhost:8080/%s", shortURL) // modify the base URL as needed

	fmt.Fprintf(w, "SHort URL: %s", html.EscapeString(shortURL))
}

func base62Encode(n *big.Int) string {
	base := big.NewInt(62)
	encoded := ""
	zero := big.NewInt(0)

	for n.Cmp(zero) > 0 {
		quotient, remainder := new(big.Int), new(big.Int)
		n.DivMod(n, base, remainder)
		encoded = string(base62Charset[remainder.Int64()]) + encoded
	}

	return encoded
}

func main() {
	http.HandleFunc("/", redirectToLongURL)
	http.HandleFunc("/shorten", shortenURL)
	http.ListenAndServe(":8080", nil)
}
