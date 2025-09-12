package mysqlrepo

import (
	"database/sql"
	"errors"
	"time"

	domuser "github.com/fintrack/user-service/internal/core/domain/entities/user"
	domerrors "github.com/fintrack/user-service/internal/core/errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(email string) (*domuser.User, error) {
	const q = `SELECT id, email, password_hash, first_name, last_name, role, is_active, email_verified, created_at, updated_at FROM users WHERE email = ? LIMIT 1`
	u := domuser.User{}
	var role string
	err := r.db.QueryRow(q, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &role, &u.IsActive, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domerrors.ErrUserNotFound
		}
		return nil, err
	}
	u.Role = domuser.Role(role)
	return &u, nil
}

func (r *UserRepository) GetByID(id string) (*domuser.User, error) {
	const q = `SELECT id, email, password_hash, first_name, last_name, role, is_active, email_verified, created_at, updated_at FROM users WHERE id = ? LIMIT 1`
	u := domuser.User{}
	var role string
	err := r.db.QueryRow(q, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &role, &u.IsActive, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domerrors.ErrUserNotFound
		}
		return nil, err
	}
	u.Role = domuser.Role(role)
	return &u, nil
}

func (r *UserRepository) Create(u *domuser.User) error {
	const q = `INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active, email_verified, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
	_, err := r.db.Exec(q, u.ID, u.Email, u.PasswordHash, u.FirstName, u.LastName, string(u.Role), u.IsActive, u.EmailVerified, u.CreatedAt, u.UpdatedAt)
	return err
}
