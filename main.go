package main

import (
    "embed"
    "fmt"
    "net/http"
    "os"

    "github.com/qazer2687/bin/internal/config"
    "github.com/qazer2687/bin/internal/handler"
)

// declare a variable "static" that holds all files under the src/static director at compile time
//
//go:embed static/*
var static embed.FS

func main() {
	// load configuration
	cfg := config.Load()
    h := handler.New(&cfg, static)
	// make sure storage path exists
	os.MkdirAll(cfg.StoragePath, 0755)

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.HandleRoot)
	// serve static assets (fonts, etc)
	mux.HandleFunc("GET /fonts/", func(w http.ResponseWriter, r *http.Request) {
		data, err := static.ReadFile("static" + r.URL.Path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "font/woff2")
		w.Write(data)
	})
	// take a file to upload to bin
	mux.HandleFunc("POST /upload", h.HandleUpload)
	// take an id to fetch a file from the server
	mux.HandleFunc("GET /{id}", h.HandleDownload)

	fmt.Printf("[INFO] server starting on %s\n", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, mux); err != nil {
		fmt.Printf("[ERROR] server failed to start: %s\n", err)
	}
}