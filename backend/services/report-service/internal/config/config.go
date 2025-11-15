package config

import (
	"os"
	"time"
)

// Config estructura de configuración del servicio
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logging  LoggingConfig
	CORS     CORSConfig
}

// ServerConfig configuración del servidor HTTP
type ServerConfig struct {
	Port         string
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig configuración de base de datos
type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// LoggingConfig configuración de logging
type LoggingConfig struct {
	Level string
}

// CORSConfig configuración de CORS
type CORSConfig struct {
	AllowedOrigins []string
}

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8085"),
			Environment:  getEnv("GIN_MODE", "debug"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "mysql"),
			Port:     getEnv("DB_PORT", "3306"),
			Name:     getEnv("DB_NAME", "fintrack"),
			User:     getEnv("DB_USER", "fintrack_user"),
			Password: getEnv("DB_PASSWORD", "fintrack_password"),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{
				getEnv("ALLOWED_ORIGINS", "http://localhost:4200"),
			},
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
