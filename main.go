package main

import (   
    "os"
    "net/http"
    "github.com/go-chi/chi/v5"
    "conecto/api"
    "conecto/pipelines"
    "github.com/joho/godotenv"
    "log"
)

func main() {

    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    r := chi.NewRouter()

    r.Handle("/*", http.FileServer(http.Dir("./clients/ms-excel")))


    r.Get("/proxy", api.ProxyHandler)

    http.ListenAndServe(":"+port, r)

    registry := pipelines.NewRegistry()

    pipeline := registry.Factories["facebookAdInsight"]()

    err := pipeline.Run(ctx)
}   