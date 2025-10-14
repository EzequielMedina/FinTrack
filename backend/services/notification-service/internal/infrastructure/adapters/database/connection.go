package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/fintrack/notification-service/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

// Connection estructura para manejar la conexión a la base de datos
type Connection struct {
	DB *sql.DB
}

// NewConnection crea una nueva conexión a la base de datos
func NewConnection(cfg *config.DatabaseConfig) (*Connection, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Configurar parámetros de conexión
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verificar la conexión
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return &Connection{DB: db}, nil
}

// Close cierra la conexión a la base de datos
func (c *Connection) Close() error {
	return c.DB.Close()
}

// CreateTables ya no es necesario - las tablas se crean via migraciones
// Ver: database/migrations/08_V8__notifications.sql
func (c *Connection) CreateTables() error {
	// Las tablas notification service se crean automáticamente via migraciones
	log.Println("✅ Tablas notification service manejadas via migraciones")
	return nil
}
