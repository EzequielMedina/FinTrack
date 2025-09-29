package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fintrack/transaction-service/internal/infrastructure/entrypoints/router"
	_ "github.com/go-sql-driver/mysql"
)

// Config holds application configuration
type Config struct {
	Port       int
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		Port:       8080, // default port
		DBHost:     "localhost",
		DBPort:     3306,
		DBUser:     "root",
		DBPassword: "password",
		DBName:     "fintrack",
	}

	// Load from environment variables
	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Port = p
		}
	}

	// Default port for transaction service if not set
	if config.Port == 8080 {
		config.Port = 8083 // Default for transaction service
	}

	if host := os.Getenv("DB_HOST"); host != "" {
		config.DBHost = host
	}

	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.DBPort = p
		}
	}

	if user := os.Getenv("DB_USER"); user != "" {
		config.DBUser = user
	}

	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.DBPassword = password
	}

	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.DBName = dbName
	}

	return config
}

// InitDatabase initializes database connection
func InitDatabase(config *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}

// CORS middleware for simple HTTP
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-User-ID")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Starting Transaction Service...")

	// Load configuration
	config := LoadConfig()
	log.Printf("Configuration loaded: Port=%d, DB=%s:%d/%s", config.Port, config.DBHost, config.DBPort, config.DBName)

	// Initialize database
	db, err := InitDatabase(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Println("Database connection established")

	// Create router
	appRouter := router.NewRouter(db)
	mux := appRouter.SetupRoutes()

	// Add CORS middleware
	handler := corsMiddleware(mux)

	// Start server
	addr := fmt.Sprintf(":%d", config.Port)
	log.Printf("Server starting on http://localhost%s", addr)
	log.Printf("Health check: http://localhost%s/health", addr)
	log.Printf("API endpoints: http://localhost%s/api/v1/transactions", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
