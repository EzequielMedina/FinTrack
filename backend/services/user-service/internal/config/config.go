package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Port          string
	DBHost        string
	DBPort        string
	DBName        string
	DBUser        string
	DBPassword    string
	JWTSecret     string
	JWTExpiry     time.Duration
	RefreshExpiry time.Duration
	LogLevel      string
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func ParseDurationEnv(key, def string) time.Duration {
	s := getenv(key, def)
	d, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:          getenv("PORT", "8080"),
		DBHost:        getenv("DB_HOST", "localhost"),
		DBPort:        getenv("DB_PORT", "3306"),
		DBName:        getenv("DB_NAME", "fintrack"),
		DBUser:        getenv("DB_USER", "fintrack_user"),
		DBPassword:    getenv("DB_PASSWORD", "fintrack_password"),
		JWTSecret:     getenv("JWT_SECRET", "change-me"),
		JWTExpiry:     ParseDurationEnv("JWT_EXPIRY", "24h"),
		RefreshExpiry: ParseDurationEnv("JWT_REFRESH_EXPIRY", "168h"),
		LogLevel:      getenv("LOG_LEVEL", "info"),
	}
	if cfg.JWTSecret == "change-me" {
		// not fatal but warn; keep simple
		fmt.Fprintf(os.Stderr, "[WARN] using default JWT secret, please set JWT_SECRET\n")
	}
	return cfg, nil
}
