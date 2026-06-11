package main

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// declare a variable "static" that holds all files under the src/static director at compile time
//
//go:embed static/*
var static embed.FS

// load config struct
var cfg = loadConfig()
var uploadSem = make(chan struct{}, cfg.MaxUploads)

func main() {
	// make sure storage path exists
	os.MkdirAll(cfg.StoragePath, 0755)

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
	data, err := static.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Write(data)
}

type throttledWriter struct {
	w    io.Writer
	rate int64
}

func (t *throttledWriter) Write(p []byte) (int, error) {
	if t.rate > 0 {
		time.Sleep(time.Duration(len(p)) * time.Second / time.Duration(t.rate))
	}
	return t.w.Write(p)
}

func handleUpload(
	w http.ResponseWriter,
	r *http.Request,
) {
	select {
	case uploadSem <- struct{}{}:
		defer func() { <-uploadSem }()
	default:
		http.Error(w, "too many uploads", http.StatusServiceUnavailable)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, cfg.MaxFileSize)

	temp, err := os.CreateTemp("", "b1n-*")
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	hash := sha256.New()
	_, err = io.Copy(&throttledWriter{temp, cfg.UploadRate}, io.TeeReader(r.Body, hash))
	temp.Close()
	if err != nil {
		os.Remove(temp.Name())
		http.Error(w, "body too large", http.StatusRequestEntityTooLarge)
		return
	}
	id := hex.EncodeToString(hash.Sum(nil))

	os.Rename(temp.Name(), filepath.Join(cfg.StoragePath, id))
	os.Remove(temp.Name())

	fmt.Fprint(w, id)
}

func handleDownload(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("id")
	// get the full path of the file on the server
	path := filepath.Join(cfg.StoragePath, id)

	file, err := os.Open(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// ensure the file is closed however the function returns
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.ServeContent(w, r, id, stat.ModTime(), file)
}
