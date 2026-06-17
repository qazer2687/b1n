package main

import (
    "embed"
    "fmt"
    "net/http"
    "os"

    "github.com/qazer2687/b1n/internal/config"
    "github.com/qazer2687/b1n/internal/handler"
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

	// set up the request router
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.HandleRoot)
	// serve static assets (fonts, images, etc)
	mux.HandleFunc("GET /fonts/", func(w http.ResponseWriter, r *http.Request) {
		data, err := static.ReadFile("static" + r.URL.Path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "font/woff2")
		w.Write(data)
	})
	mux.HandleFunc("GET /assets/", func(w http.ResponseWriter, r *http.Request) {
		data, err := static.ReadFile("static" + r.URL.Path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if len(r.URL.Path) > 4 && r.URL.Path[len(r.URL.Path)-4:] == ".svg" {
			w.Header().Set("Content-Type", "image/svg+xml")
		}
		w.Write(data)
	})
	
	// serve the stylesheet
	mux.HandleFunc("GET /style.css", func(w http.ResponseWriter, r *http.Request) {
	    data, err := static.ReadFile("static/style.css")
	    if err != nil {
	        http.NotFound(w, r)
	        return
	    }
	    w.Header().Set("Content-Type", "text/css")
	    w.Write(data)
	})
	// serve the client-side script
	mux.HandleFunc("GET /script.js", func(w http.ResponseWriter, r *http.Request) {
	    data, err := static.ReadFile("static/script.js")
	    if err != nil {
	        http.NotFound(w, r)
	        return
	    }
	    w.Header().Set("Content-Type", "application/javascript")
	    w.Write(data)
	})
	
	// take a file to upload to b1n
	mux.HandleFunc("POST /upload", h.HandleUpload)
	// take an id to fetch a file from the server
	mux.HandleFunc("GET /{id}", h.HandleDownload)

	// start the http server
	fmt.Printf("[INFO] server starting on %s\n", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, mux); err != nil {
		fmt.Printf("[ERROR] server failed to start: %s\n", err)
	}
}