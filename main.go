package main

import (
	"math/big"
	"net/http"
	"sync"
)

var (
	shortURLS     = make(map[string]string) // Map to store short URLs and their corresponding long URLs
	urlCounter    = new(big.Int)            // Coubnter to generate unique URL IDs
	urlLock       = sync.RWMutex{}          // // Mutex to synchronize access to the shortURLs map
	base62Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func redirectToLongURL(w http.ResponseWriter, r *http.Request) {
	short := r.URL.Path[1:] // remove the leading slash from the URL

	//check is the chort URL exists in the map
	urlLock.RLock()
	longURL, exists := shortURLS[short]
	urlLock.Unlock()

	if exists {
		http.Redirect(w, r, longURL, http.StatusMovedPermanently) // redirect to long URL
	}
	if !exists {
		http.NotFound(w, r)
	}
}
