package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         string // listen address
	StoragePath  string // directory for uploaded files
	MaxFileSize  int64  // per-file upload limit
	MaxTotalSize int64  // total storage cap
	BaseURL      string // public-facing url for generated links
	MaxUploads   int    // max concurrent uploads
	UploadRate   int64  // bytes per second upload throttle
}

func Load() Config {
	return Config{
		Port:         env("PORT", ":8080"),
		StoragePath:  env("STORAGE_PATH", "./data"),
		MaxFileSize:  envInt("MAX_FILE_SIZE", 10737418240),
		MaxTotalSize: envInt("MAX_TOTAL_SIZE", 10737418240),
		BaseURL:      env("BASE_URL", "http://localhost:8080"),
		MaxUploads:   int(envInt("MAX_UPLOADS", 1)),
		UploadRate:   envInt("UPLOAD_RATE", 0),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return fallback
	}
	return n
}
