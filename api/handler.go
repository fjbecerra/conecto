package api

import (
	"conecto/pipelines"
	"conecto/sinks"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
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

func (r *pipelines.Registry) Handle(w http.ResponseWriter, req *http.Request) {
    source := req.URL.Query().Get("source")

    factory, ok := r.Factories[source]
    if !ok {
        http.Error(w, "unknown source", 400)
        return
    }

    ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
    defer cancel()

    pipeline := factory()
    sink := sinks.NewMemorySink()

    records, _ := pipeline.Connector.Run(ctx)
    extracted, _ := pipeline.Extractor.Extract(ctx, records)
    computed, _ := pipeline.Computer.Compute(ctx, extracted)

    if err := sink.Write(ctx, computed); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    json.NewEncoder(w).Encode(sink.GetData())
}