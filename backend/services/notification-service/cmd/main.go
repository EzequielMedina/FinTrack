package main

import (
	"log"
	"os"

	"github.com/fintrack/notification-service/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde .env en desarrollo
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	log.Println("🚀 Iniciando FinTrack Notification Service...")

	// Crear e inicializar la aplicación
	application := app.NewApplication()

	// Iniciar el servidor
	if err := application.Start(); err != nil {
		log.Printf("❌ Error iniciando el servidor: %v", err)
		os.Exit(1)
	}
}
