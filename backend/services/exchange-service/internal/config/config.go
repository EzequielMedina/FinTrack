package config

import (
	"os"
	"strconv"
	"time"
)

// Config estructura de configuración del servicio
type Config struct {
	Server   ServerConfig
	DolarAPI DolarAPIConfig
	Logging  LoggingConfig
}

// ServerConfig configuración del servidor HTTP
type ServerConfig struct {
	Port         string
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DolarAPIConfig configuración de la API de DolarAPI
type DolarAPIConfig struct {
	BaseURL string
	Timeout time.Duration
}

// LoggingConfig configuración de logging
type LoggingConfig struct {
	Level string
}

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8087"),
			Environment:  getEnv("GIN_MODE", "debug"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 30*time.Second),
		},
		DolarAPI: DolarAPIConfig{
			BaseURL: getEnv("DOLAR_API_BASE_URL", "https://dolarapi.com"),
			Timeout: getDurationEnv("DOLAR_API_TIMEOUT", 10*time.Second),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}
}

// getEnv obtiene una variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDurationEnv obtiene una duración desde variable de entorno
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getIntEnv obtiene un entero desde variable de entorno
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
