package handler

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

	"github.com/qazer2687/bin/internal/config"
)

type Handler struct {
	cfg       *config.Config
	uploadSem chan struct{}
	static    embed.FS
}

func New(cfg *config.Config, static embed.FS) *Handler {
	return &Handler{
		cfg:       cfg,
		uploadSem: make(chan struct{}, cfg.MaxUploads),
		static:    static,
	}
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

func (h *Handler) HandleRoot(
	w http.ResponseWriter,
	r *http.Request,
) {
	data, err := h.static.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Write(data)
}

func (h *Handler) HandleUpload(
	w http.ResponseWriter,
	r *http.Request,
) {
	select {
	case h.uploadSem <- struct{}{}:
		defer func() { <-h.uploadSem }()
	default:
		http.Error(w, "too many uploads", http.StatusServiceUnavailable)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, h.cfg.MaxFileSize)
	temp, err := os.CreateTemp(h.cfg.StoragePath, "_*")
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	hash := sha256.New()
	_, err = io.Copy(&throttledWriter{temp, h.cfg.UploadRate}, io.TeeReader(r.Body, hash))
	temp.Close()
	if err != nil {
		os.Remove(temp.Name())
		http.Error(w, "body too large", http.StatusRequestEntityTooLarge)
		return
	}
	id := hex.EncodeToString(hash.Sum(nil))
	os.Rename(temp.Name(), filepath.Join(h.cfg.StoragePath, id))
	fmt.Fprint(w, id)
}

func (h *Handler) HandleDownload(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("id")
	// ignore temp files which are all prefixed with an underscore
	if len(id) > 0 && id[0] == '_' {
		http.NotFound(w, r)
		return
	}
	// get the full path of the file on the server
	path := filepath.Join(h.cfg.StoragePath, id)
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
