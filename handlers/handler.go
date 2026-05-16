package api

import (
	
	"io"
	"net/http"
	"net/url"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
    // Get the target URL from query param ?url=https://api.example.com
    encodedUrl := r.URL.Query().Get("url")
    if encodedUrl == "" {
        http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
        return
    }

    // Decode it
	decodedUrl, err := url.QueryUnescape(encodedUrl)
	if err != nil {
		http.Error(w, "Invalid URL encoding", http.StatusBadRequest)
		return
	}

    // Make the HTTP request
    resp, err := http.Get(decodedUrl)
    if err != nil {
        http.Error(w, "Failed to fetch URL: "+err.Error(), http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    // Copy the response headers
    for key, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }

    // Set the status code
    w.WriteHeader(resp.StatusCode)

    // Copy the body
    io.Copy(w, resp.Body)
}
