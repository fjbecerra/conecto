package api

import (
    "io"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestProxyHandler(t *testing.T) {
    // Step 1: Create a fake server to act as the external API
    fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain")
        w.WriteHeader(http.StatusOK)
        io.WriteString(w, "Hello from fake API")
    }))
    defer fakeServer.Close()

    // Step 2: Create a request to your ProxyHandler with the fake URL
    req := httptest.NewRequest(http.MethodGet, "/proxy?url="+fakeServer.URL, nil)
    w := httptest.NewRecorder()

    // Step 3: Call the handler
    ProxyHandler(w, req)

    // Step 4: Check the response
    res := w.Result()
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        t.Errorf("expected status 200, got %d", res.StatusCode)
    }

    body, err := io.ReadAll(res.Body)
    if err != nil {
        t.Fatalf("could not read response: %v", err)
    }

    expected := "Hello from fake API"
    if strings.TrimSpace(string(body)) != expected {
        t.Errorf("expected body %q, got %q", expected, string(body))
    }
}