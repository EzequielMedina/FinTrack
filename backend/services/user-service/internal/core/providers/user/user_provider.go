package userprovider

import (
	domuser "github.com/fintrack/user-service/internal/core/domain/entities/user"
)

type Repository interface {
	GetByEmail(email string) (*domuser.User, error)
	GetByID(id string) (*domuser.User, error)
	Create(u *domuser.User) error
}
