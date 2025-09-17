package router

import (
	"github.com/fintrack/account-service/internal/app"
	accounthandler "github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/account"
)

type Handlers struct {
	Account *accounthandler.Handler
}

func NewHandlers(a *app.Application) *Handlers {
	return &Handlers{
		Account: accounthandler.New(a.AccountService),
	}
}
