package app

import (
	"fmt"
	"time"

	"github.com/fintrack/account-service/internal/config"
	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/core/service"
	mysqlrepo "github.com/fintrack/account-service/internal/infrastructure/repositories/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Application struct {
	Config             *config.Config
	DB                 *gorm.DB
	AccountService     *service.AccountService
	CardService        *service.CardService
	InstallmentService *service.InstallmentService
}

func New(cfg *config.Config) (*Application, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Auto-migrate tables
	if err := gormDB.AutoMigrate(&entities.Account{}, &entities.Card{}, &entities.InstallmentPlan{}, &entities.Installment{}); err != nil {
		return nil, fmt.Errorf("failed to migrate tables: %w", err)
	}

	// repositories
	accountRepo := mysqlrepo.NewAccountRepository(gormDB)
	cardRepo := mysqlrepo.NewCardRepository(gormDB)
	installmentRepo := mysqlrepo.NewInstallmentRepository(gormDB)
	installmentPlanRepo := mysqlrepo.NewInstallmentPlanRepository(gormDB)

	// services
	accountSvc := service.NewAccountService(accountRepo)
	installmentSvc := service.NewInstallmentService(installmentRepo, installmentPlanRepo, cardRepo)
	cardSvc := service.NewCardService(cardRepo, accountRepo, installmentSvc)

	return &Application{
		Config:             cfg,
		DB:                 gormDB,
		AccountService:     accountSvc,
		CardService:        cardSvc,
		InstallmentService: installmentSvc,
	}, nil
}

func (a *Application) Close() error {
	if a.DB != nil {
		sqlDB, err := a.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
