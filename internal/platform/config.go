package platform

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv        string
	AppPort       string
	DatabaseURL   string
	JWTSecret     string
	JWTExpiryHour time.Duration
}

func Load() *Config {
	loadEnvFile()

	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL == "" {
		panic("DATABASE_URL is not set; check your .env file and working directory")
	}

	return &Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		AppPort:       getEnv("APP_PORT", "8080"),
		DatabaseURL:   databaseURL,
		JWTSecret:     getEnv("JWT_SECRET", "dev-secret-change-me"),
		JWTExpiryHour: getDurationEnv("JWT_EXPIRY_HOURS", 24),
	}
}

func loadEnvFile() {
	// Try current directory first
	if err := godotenv.Load(); err == nil {
		return
	}

	// Try parent directory (when running from ./tmp via air)
	if err := godotenv.Load("../.env"); err == nil {
		return
	}

	// Try to find project root by looking for go.mod
	if root := findProjectRoot(); root != "" {
		_ = godotenv.Load(filepath.Join(root, ".env"))
	}
}

func findProjectRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDurationEnv(key string, fallbackHours int) time.Duration {
	if v := os.Getenv(key); v != "" {
		if hours, err := strconv.Atoi(v); err == nil {
			return time.Duration(hours) * time.Hour
		}
	}
	return time.Duration(fallbackHours) * time.Hour
}
