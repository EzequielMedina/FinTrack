package app

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/fintrack/user-service/internal/config"
	"github.com/fintrack/user-service/internal/core/service"
	mysqlrepo "github.com/fintrack/user-service/internal/infrastructure/repositories/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type Application struct {
	Config      *config.Config
	DB          *sql.DB
	AuthService *service.AuthService
}

func New(cfg *config.Config) (*Application, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// repositories
	userRepo := mysqlrepo.NewUserRepository(db)

	// services
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry, cfg.RefreshExpiry)

	return &Application{
		Config:      cfg,
		DB:          db,
		AuthService: authSvc,
	}, nil
}
