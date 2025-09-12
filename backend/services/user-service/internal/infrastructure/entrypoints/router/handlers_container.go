package router

import (
	"github.com/fintrack/user-service/internal/app"
	authhandler "github.com/fintrack/user-service/internal/infrastructure/entrypoints/handlers/auth"
)

type Handlers struct {
	Auth *authhandler.Handler
}

func NewHandlers(a *app.Application) *Handlers {
	return &Handlers{
		Auth: authhandler.New(a.AuthService),
	}
}
