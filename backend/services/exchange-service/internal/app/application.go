package app

import (
	"log"

	"github.com/fintrack/exchange-service/internal/config"
	"github.com/fintrack/exchange-service/internal/core/service"
	"github.com/fintrack/exchange-service/internal/infrastructure/clients"
	exchangehandler "github.com/fintrack/exchange-service/internal/infrastructure/entrypoints/handlers/exchange"
	"github.com/fintrack/exchange-service/internal/infrastructure/entrypoints/router"
	"github.com/fintrack/exchange-service/internal/infrastructure/providers"
	"github.com/gin-gonic/gin"
)

// Application contiene todas las dependencias del servicio
type Application struct {
	Config          *config.Config
	ExchangeService *service.ExchangeService
	Router          *gin.Engine
}

// NewApplication crea e inicializa la aplicación con todas sus dependencias
func NewApplication() *Application {
	log.Println("Inicializando aplicación Exchange Service...")

	// Cargar configuración
	cfg := config.LoadConfig()
	log.Printf("Configuración cargada - Puerto: %s, Environment: %s", cfg.Server.Port, cfg.Server.Environment)

	// Configurar modo Gin
	gin.SetMode(cfg.Server.Environment)

	// Crear cliente HTTP para DolarAPI
	dolarAPIClient := clients.NewDolarAPIClient(cfg.DolarAPI.BaseURL, cfg.DolarAPI.Timeout)
	log.Printf("Cliente DolarAPI creado - Base URL: %s", cfg.DolarAPI.BaseURL)

	// Crear proveedor de exchange rates
	exchangeProvider := providers.NewDolarAPIProvider(dolarAPIClient)
	log.Println("Proveedor de exchange rates creado")

	// Crear servicio de exchange
	exchangeService := service.NewExchangeService(exchangeProvider)
	log.Println("Servicio de exchange creado")

	// Crear handler
	exchangeHandler := exchangehandler.New(exchangeService)
	log.Println("Handler de exchange creado")

	// Configurar rutas
	router := router.SetupRoutes(exchangeHandler)
	log.Println("Router configurado")

	log.Println("Aplicación inicializada exitosamente")

	return &Application{
		Config:          cfg,
		ExchangeService: exchangeService,
		Router:          router,
	}
}

// Start inicia el servidor HTTP
func (app *Application) Start() error {
	log.Printf("Iniciando servidor en puerto %s...", app.Config.Server.Port)

	return app.Router.Run(":" + app.Config.Server.Port)
}
