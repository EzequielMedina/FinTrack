package config

import (
	"os"
)

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type LLMConfig struct {
	Provider string // "groq" o "ollama"
	// Groq config
	GroqAPIKey string
	GroqModel  string
	// Ollama config
	OllamaHost    string
	OllamaModel   string
	OllamaTimeout string
}

type ReportConfig struct {
	Engine string
}

type AppConfig struct {
	Server   ServerConfig
	Database DatabaseConfig
	LLM      LLMConfig
	Report   ReportConfig
	Timezone string
}

func env(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func Load() *AppConfig {
	return &AppConfig{
		Server: ServerConfig{Port: env("PORT", "8090")},
		Database: DatabaseConfig{
			Host:     env("DB_HOST", "mysql"),
			Port:     env("DB_PORT", "3306"),
			Name:     env("DB_NAME", "fintrack"),
			User:     env("DB_USER", "fintrack_user"),
			Password: env("DB_PASSWORD", "fintrack_password"),
		},
		LLM: LLMConfig{
			Provider: env("LLM_PROVIDER", "groq"), // Default a Groq
			// Groq config
			GroqAPIKey: env("GROQ_API_KEY", ""),
			GroqModel:  env("GROQ_MODEL", "llama-3.1-8b-instant"),
			// Ollama config (fallback)
			OllamaHost:    env("OLLAMA_HOST", "http://localhost:11434"),
			OllamaModel:   env("OLLAMA_MODEL", "qwen2.5:3b"),
			OllamaTimeout: env("OLLAMA_TIMEOUT", "30s"),
		},
		Report:   ReportConfig{Engine: env("REPORT_PDF_ENGINE", "gofpdf")},
		Timezone: env("TIMEZONE", "America/Argentina/Buenos_Aires"),
	}
}
