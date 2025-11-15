package app

import (
	"log"

	"github.com/fintrack/chatbot-service/internal/config"
	"github.com/fintrack/chatbot-service/internal/core/ports"
	"github.com/fintrack/chatbot-service/internal/core/service"
	"github.com/fintrack/chatbot-service/internal/infrastructure/entrypoints/handlers"
	"github.com/fintrack/chatbot-service/internal/infrastructure/router"
	"github.com/fintrack/chatbot-service/internal/providers/data/mysql"
	"github.com/fintrack/chatbot-service/internal/providers/groq"
	"github.com/fintrack/chatbot-service/internal/providers/ollama"
	"github.com/fintrack/chatbot-service/internal/providers/pdf"
	"github.com/gin-gonic/gin"
)

type Application struct {
	Config *config.AppConfig
	Router *gin.Engine
}

func Initialize(cfg *config.AppConfig) (*Application, error) {
	// DB
	conn, err := mysql.NewConnection(cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Database.User, cfg.Database.Password)
	if err != nil {
		return nil, err
	}
	dataProv := mysql.NewDataProvider(conn)
	log.Println("✅ Conexión MySQL establecida")

	// LLM Provider (Groq o Ollama)
	var llmClient ports.LLMProvider
	switch cfg.LLM.Provider {
	case "groq":
		if cfg.LLM.GroqAPIKey == "" {
			log.Printf("⚠️  GROQ_API_KEY no configurada, usando Ollama como fallback")
			llmClient = ollama.NewClient(cfg.LLM.OllamaHost, cfg.LLM.OllamaModel)
			log.Printf("🦙 Usando Ollama LLM fallback en %s con modelo %s", cfg.LLM.OllamaHost, cfg.LLM.OllamaModel)
		} else {
			llmClient = groq.NewGroqClient(cfg.LLM.GroqAPIKey, cfg.LLM.GroqModel)
			log.Printf("🚀 Usando Groq LLM con modelo %s", cfg.LLM.GroqModel)
		}
	case "ollama":
		llmClient = ollama.NewClient(cfg.LLM.OllamaHost, cfg.LLM.OllamaModel)
		log.Printf("🦙 Usando Ollama LLM en %s con modelo %s", cfg.LLM.OllamaHost, cfg.LLM.OllamaModel)
	default:
		log.Printf("⚠️  LLM_PROVIDER no válido: %s, usando Ollama como fallback", cfg.LLM.Provider)
		llmClient = ollama.NewClient(cfg.LLM.OllamaHost, cfg.LLM.OllamaModel)
	}

	pdfGen := pdf.NewGenerator()

	// Service
	chatbotSvc := service.NewChatbotService(dataProv, llmClient, pdfGen)
	handler := handlers.NewChatHandler(chatbotSvc)

	// Router
	r := router.SetupRoutes(handler)
	log.Println("🛤️  Router configurado en chatbot-service")

	return &Application{Config: cfg, Router: r}, nil
}

func (a *Application) Start() error {
	log.Printf("🚀 Iniciando Chatbot Service en puerto %s...", a.Config.Server.Port)
	return a.Router.Run(":" + a.Config.Server.Port)
}
