package main

import (
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"sync"
	"strings"
)

func main() {
	// load config struct
	cfg := loadConfig();

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	// take a file to upload to bin
	mux.HandleFunc("POST /upload", handleUpload)
	// take an id to fetch a file from the server
	mux.HandleFunc("GET /{id}", handleDownload)

	if err := http.ListenAndServe(cfg.Port, mux); err != nil {
	    fmt.Printf("[ERROR] server failed to start: %s\n", err)
	}
	fmt.Printf("[INFO] server started on port %s\n", cfg.Port)

	
}

func handleRoot(
	w http.ResponseWriter,
	r *http.Request,
) {
	fmt.Fprintf(w, "ok")
}

func handleUpload(
	w http.ResponseWriter,
	r *http.Request,
) {
	
}

func handleDownload(
	w http.ResponseWriter,
	r *http.Request,
) {
	
}

func hash(s string) uint32 {
	hash := fnv.New32()
	hash.Write([]byte(s))
	return hash.Sum32()
}
