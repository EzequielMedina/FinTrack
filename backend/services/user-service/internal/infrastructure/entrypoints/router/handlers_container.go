package router

import (
	"github.com/fintrack/user-service/internal/app"
	authhandler "github.com/fintrack/user-service/internal/infrastructure/entrypoints/handlers/auth"
	userhandler "github.com/fintrack/user-service/internal/infrastructure/entrypoints/handlers/user"
)

type Handlers struct {
	Auth *authhandler.Handler
	User *userhandler.Handler
}

func NewHandlers(a *app.Application) *Handlers {
	return &Handlers{
		Auth: authhandler.New(a.AuthService),
		User: userhandler.New(a.UserService),
	}
}
