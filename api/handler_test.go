package api

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "encoding/json"
)

func TestHelloHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/hello", nil)
    w := httptest.NewRecorder()

    HelloHandler(w, req)

    res := w.Result()
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        t.Errorf("expected status 200, got %d", res.StatusCode)
    }

    var msg Message
    if err := json.NewDecoder(res.Body).Decode(&msg); err != nil {
        t.Fatalf("could not decode response: %v", err)
    }

    expected := "Hello from API folder"
    if msg.Text != expected {
        t.Errorf("expected message %q, got %q", expected, msg.Text)
    }
}