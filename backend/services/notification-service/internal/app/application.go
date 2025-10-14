package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fintrack/notification-service/internal/config"
	"github.com/fintrack/notification-service/internal/core/service"
	"github.com/fintrack/notification-service/internal/infrastructure/adapters/database"
	"github.com/fintrack/notification-service/internal/infrastructure/adapters/email"
	"github.com/fintrack/notification-service/internal/infrastructure/entrypoints/handlers/notification"
	"github.com/fintrack/notification-service/internal/infrastructure/entrypoints/router"
	"github.com/fintrack/notification-service/internal/infrastructure/jobs"
	"github.com/gin-gonic/gin"
)

// Application contiene todas las dependencias del servicio
type Application struct {
	Config              *config.Config
	NotificationService *service.NotificationService
	JobScheduler        *jobs.JobScheduler
	Router              *gin.Engine
	DBConnection        *database.Connection
}

// NewApplication crea e inicializa la aplicación con todas sus dependencias
func NewApplication() *Application {
	log.Println("🚀 Inicializando Notification Service...")

	// Cargar configuración
	cfg := config.LoadConfig()
	log.Printf("📋 Configuración cargada - Puerto: %s, Environment: %s", cfg.Server.Port, cfg.Server.Environment)

	// Configurar modo Gin
	gin.SetMode(cfg.Server.Environment)

	// Conectar a la base de datos
	log.Printf("🔌 Conectando a la base de datos: %s:%s/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	dbConnection, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("❌ Error conectando a la base de datos: %v", err)
	}
	log.Println("✅ Conexión a base de datos establecida")

	// Crear tablas necesarias
	if err := dbConnection.CreateTables(); err != nil {
		log.Fatalf("❌ Error creando tablas: %v", err)
	}
	log.Println("✅ Tablas de base de datos verificadas")

	// Crear repositorios
	cardRepo := database.NewCardRepository(dbConnection.DB)
	installmentRepo := database.NewInstallmentRepository(dbConnection.DB)
	notificationRepo := database.NewNotificationRepository(dbConnection.DB)
	log.Println("✅ Repositorios creados")

	// Crear cliente EmailJS
	emailClient := email.NewEmailJSClient(&cfg.EmailJS)
	log.Printf("📧 Cliente EmailJS creado - Service ID: %s", cfg.EmailJS.ServiceID)

	// Crear servicio de notificaciones
	notificationService := service.NewNotificationService(
		cardRepo,
		installmentRepo,
		notificationRepo,
		emailClient,
	)
	log.Println("✅ Servicio de notificaciones creado")

	// Crear job scheduler
	jobScheduler := jobs.NewJobScheduler(notificationService, &cfg.Job)
	log.Println("⏰ Job scheduler creado")

	// Crear handler
	notificationHandler := notification.New(notificationService, jobScheduler)
	log.Println("✅ Handler de notificaciones creado")

	// Configurar rutas
	routerGin := router.SetupRoutes(notificationHandler)
	log.Println("🛤️  Router configurado")

	log.Println("🎉 Aplicación inicializada exitosamente")

	return &Application{
		Config:              cfg,
		NotificationService: notificationService,
		JobScheduler:        jobScheduler,
		Router:              routerGin,
		DBConnection:        dbConnection,
	}
}

// Start inicia el servidor HTTP y el job scheduler
func (app *Application) Start() error {
	log.Printf("🚀 Iniciando Notification Service en puerto %s...", app.Config.Server.Port)

	// Iniciar job scheduler
	if err := app.JobScheduler.Start(); err != nil {
		log.Printf("⚠️  Error iniciando job scheduler: %v", err)
		return err
	}

	// Canal para manejar señales del sistema
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor HTTP en una goroutine
	go func() {
		log.Printf("🌐 Servidor HTTP iniciado en puerto %s", app.Config.Server.Port)
		if err := app.Router.Run(":" + app.Config.Server.Port); err != nil {
			log.Printf("❌ Error del servidor HTTP: %v", err)
		}
	}()

	// Esperar señal de terminación
	sig := <-signalChan
	log.Printf("📞 Recibida señal: %v", sig)

	// Graceful shutdown
	return app.Shutdown()
}

// Shutdown realiza un cierre limpio del servicio
func (app *Application) Shutdown() error {
	log.Println("🛑 Iniciando cierre limpio del servicio...")

	// Detener job scheduler
	app.JobScheduler.Stop()

	// Cerrar conexión a base de datos
	if err := app.DBConnection.Close(); err != nil {
		log.Printf("⚠️  Error cerrando conexión a base de datos: %v", err)
	} else {
		log.Println("✅ Conexión a base de datos cerrada")
	}

	log.Println("✅ Servicio cerrado exitosamente")
	return nil
}
