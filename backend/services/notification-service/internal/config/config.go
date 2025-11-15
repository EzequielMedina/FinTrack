package config

import (
	"os"
	"strconv"
	"time"
)

// Config estructura de configuración del servicio
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	EmailJS  EmailJSConfig
	Job      JobConfig
	Logging  LoggingConfig
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

// EmailJSConfig configuración de EmailJS API
type EmailJSConfig struct {
	ServiceID  string
	TemplateID string
	PublicKey  string
	PrivateKey string
	FromName   string
	ReplyTo    string
}

// JobConfig configuración del job scheduler
type JobConfig struct {
	Enabled  bool
	Schedule string
	Timezone string
}

// LoggingConfig configuración de logging
type LoggingConfig struct {
	Level string
}

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8088"),
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
		EmailJS: EmailJSConfig{
			ServiceID:  getEnv("EMAILJS_SERVICE_ID", "service_ceg7xlp"),
			TemplateID: getEnv("EMAILJS_TEMPLATE_ID", "template_e43va39"),
			PublicKey:  getEnv("EMAILJS_PUBLIC_KEY", "MSBb87-PQcXWr1gWK"),
			PrivateKey: getEnv("EMAILJS_PRIVATE_KEY", "sXLmpEZ8y2EYtCDtN5gZv"),
			FromName:   getEnv("EMAILJS_FROM_NAME", "FinTrack"),
			ReplyTo:    getEnv("EMAILJS_REPLY_TO", "support@fintrack.com"),
		},
		Job: JobConfig{
			Enabled:  getBoolEnv("JOB_ENABLED", true),
			Schedule: getEnv("JOB_SCHEDULE", "0 8 * * *"), // 8:00 AM daily
			Timezone: getEnv("JOB_TIMEZONE", "America/Argentina/Buenos_Aires"),
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

// getBoolEnv obtiene un booleano desde variable de entorno
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
