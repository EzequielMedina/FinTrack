package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fintrack/user-service/internal/app"
	"github.com/fintrack/user-service/internal/config"
	apirouter "github.com/fintrack/user-service/internal/infrastructure/entrypoints/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// Optional CLI health check mode for Dockerfile HEALTHCHECK compatibility
	if len(os.Args) > 1 && os.Args[1] == "health" {
		// In a more advanced version we could try pinging the DB, but here just return 0
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
	defer application.DB.Close()

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

	log.Printf("user-service listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
