package api

import (
    "encoding/json"
    "net/http"
)

type Message struct {
    Text string `json:"text"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
    response := Message{Text: "Hello from API folder"}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}