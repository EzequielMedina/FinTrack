package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/fintrack/account-service/docs"
	"github.com/fintrack/account-service/internal/app"
	"github.com/fintrack/account-service/internal/config"
	apirouter "github.com/fintrack/account-service/internal/infrastructure/entrypoints/router"
	"github.com/gin-gonic/gin"
)

// @title FinTrack Account Service API
// @version 1.0
// @description Account management service for FinTrack application
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.fintrack.com/support
// @contact.email support@fintrack.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8082
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @type apiKey

func main() {
	// Optional CLI health check mode for Dockerfile HEALTHCHECK compatibility
	if len(os.Args) > 1 && os.Args[1] == "health" {
		os.Exit(0)
		return
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
	defer application.Close()

	// Gin setup
	if cfg.LogLevel == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())

	// Basic request logger
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		log.Printf("%s %s %d %s", c.Request.Method, c.Request.URL.Path, status, latency)
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "account-service",
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	handlers := apirouter.NewHandlers(application)
	apirouter.MapRoutes(r, handlers, cfg, application)

	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("account-service listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
