package api

import (
    "io"
    "net/http"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
    // Get the target URL from query param ?url=https://api.example.com
    target := r.URL.Query().Get("url")
    if target == "" {
        http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
        return
    }

    // Make the HTTP request
    resp, err := http.Get(target)
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