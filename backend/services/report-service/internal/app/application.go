package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fintrack/report-service/internal/config"
	"github.com/fintrack/report-service/internal/core/service"
	"github.com/fintrack/report-service/internal/infrastructure/adapters/database"
	"github.com/fintrack/report-service/internal/infrastructure/entrypoints/router"
)

// Application estructura principal de la aplicación
type Application struct {
	config        *config.Config
	server        *http.Server
	reportService service.ReportService
}

// NewApplication crea una nueva instancia de la aplicación
func NewApplication() *Application {
	// Cargar configuración
	cfg := config.LoadConfig()

	// Conectar a la base de datos
	db, err := database.NewConnection(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.User,
		cfg.Database.Password,
	)
	if err != nil {
		log.Fatalf("❌ Error conectando a la base de datos: %v", err)
	}

	log.Println("✅ Conexión a base de datos establecida")

	// Crear repositorios
	reportRepo := database.NewReportRepository(db)

	// Crear servicios
	reportService := service.NewReportService(reportRepo)

	// Configurar router
	r := router.SetupRouter(cfg, reportService)

	// Crear servidor HTTP
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &Application{
		config:        cfg,
		server:        server,
		reportService: reportService,
	}
}

// Start inicia el servidor HTTP
func (a *Application) Start() error {
	// Canal para señales de interrupción
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine para iniciar el servidor
	go func() {
		log.Printf("🚀 Servidor iniciado en puerto %s", a.config.Server.Port)
		log.Printf("📊 Report Service API disponible en http://localhost:%s/api/v1", a.config.Server.Port)

		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Error en el servidor: %v", err)
		}
	}()

	// Esperar señal de interrupción
	<-quit
	log.Println("🛑 Apagando servidor...")

	// Contexto con timeout para el apagado
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Apagar el servidor gracefully
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error en el apagado del servidor: %w", err)
	}

	log.Println("✅ Servidor apagado correctamente")
	return nil
}
