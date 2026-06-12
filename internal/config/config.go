package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         string
	StoragePath  string
	MaxFileSize  int64
	MaxTotalSize int64
	BaseURL      string
	MaxUploads   int
	UploadRate   int64
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
