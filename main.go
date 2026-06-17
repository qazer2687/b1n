package main

import (
    "embed"
    "fmt"
    "io/fs"
    "mime"
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
	// register MIME types for embedded static files
	// (Alpine Linux lacks /etc/mime.types, so http.FileServer
	//  can't derive these automatically)
	mime.AddExtensionType(".svg", "image/svg+xml")
	mime.AddExtensionType(".woff2", "font/woff2")

	// load configuration
	cfg := config.Load()
    h := handler.New(&cfg, static)
	// make sure storage path exists
	os.MkdirAll(cfg.StoragePath, 0755)

	// set up the request router
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.HandleRoot)
	// serve static assets (fonts, images, etc)
	sub, _ := fs.Sub(static, "static")
	fileServer := http.FileServer(http.FS(sub))
	mux.Handle("GET /fonts/", fileServer)
	mux.Handle("GET /assets/", fileServer)
	mux.Handle("GET /style.css", fileServer)
	mux.Handle("GET /script.js", fileServer)
	
	// take a file to upload to b1n
	mux.HandleFunc("POST /upload", h.HandleUpload)
	// take an id to fetch a file from the server
	mux.HandleFunc("GET /{id}", h.HandleDownload)
	// serve raw file content for embed video playback
	mux.HandleFunc("GET /{id}/raw", h.HandleRaw)

	// start the http server
	fmt.Printf("[INFO] server starting on %s\n", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, mux); err != nil {
		fmt.Printf("[ERROR] server failed to start: %s\n", err)
	}
}