package router

import (
	"github.com/fintrack/account-service/internal/app"
	accounthandler "github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/account"
	cardhandler "github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/card"
)

type Handlers struct {
	Account *accounthandler.Handler
	Card    *cardhandler.Handler
}

func NewHandlers(a *app.Application) *Handlers {
	return &Handlers{
		Account: accounthandler.New(a.AccountService),
		Card:    cardhandler.New(a.CardService),
	}
}
