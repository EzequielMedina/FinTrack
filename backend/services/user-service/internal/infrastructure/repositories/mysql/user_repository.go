package mysqlrepo

import (
	"database/sql"
	"encoding/json"
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
	const q = `SELECT id, email, password_hash, first_name, last_name, role, is_active, email_verified, 
	           profile_data, created_at, updated_at, last_login_at FROM users WHERE email = ? LIMIT 1`
	u := domuser.User{}
	var role string
	var profileData sql.NullString
	var lastLoginAt sql.NullTime

	err := r.db.QueryRow(q, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
		&role, &u.IsActive, &u.EmailVerified, &profileData, &u.CreatedAt, &u.UpdatedAt, &lastLoginAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domerrors.ErrUserNotFound
		}
		return nil, err
	}

	u.Role = domuser.Role(role)
	if lastLoginAt.Valid {
		u.LastLoginAt = &lastLoginAt.Time
	}

	if profileData.Valid {
		if err := json.Unmarshal([]byte(profileData.String), &u.Profile); err != nil {
			// If profile data is corrupted, we'll use empty profile instead of failing
			u.Profile = domuser.Profile{}
		}
	}

	return &u, nil
}

func (r *UserRepository) GetByID(id string) (*domuser.User, error) {
	const q = `SELECT id, email, password_hash, first_name, last_name, role, is_active, email_verified, 
	           profile_data, created_at, updated_at, last_login_at FROM users WHERE id = ? LIMIT 1`
	u := domuser.User{}
	var role string
	var profileData sql.NullString
	var lastLoginAt sql.NullTime

	err := r.db.QueryRow(q, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
		&role, &u.IsActive, &u.EmailVerified, &profileData, &u.CreatedAt, &u.UpdatedAt, &lastLoginAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domerrors.ErrUserNotFound
		}
		return nil, err
	}

	u.Role = domuser.Role(role)
	if lastLoginAt.Valid {
		u.LastLoginAt = &lastLoginAt.Time
	}

	if profileData.Valid {
		if err := json.Unmarshal([]byte(profileData.String), &u.Profile); err != nil {
			u.Profile = domuser.Profile{}
		}
	}

	return &u, nil
}

func (r *UserRepository) Create(u *domuser.User) error {
	const q = `INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active, 
	           email_verified, profile_data, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now

	profileData, err := json.Marshal(u.Profile)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(q, u.ID, u.Email, u.PasswordHash, u.FirstName, u.LastName,
		string(u.Role), u.IsActive, u.EmailVerified, string(profileData), u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *UserRepository) Update(u *domuser.User) error {
	const q = `UPDATE users SET email = ?, first_name = ?, last_name = ?, role = ?, is_active = ?, 
	           email_verified = ?, profile_data = ?, updated_at = ? WHERE id = ?`

	u.UpdatedAt = time.Now()

	profileData, err := json.Marshal(u.Profile)
	if err != nil {
		return err
	}

	result, err := r.db.Exec(q, u.Email, u.FirstName, u.LastName, string(u.Role),
		u.IsActive, u.EmailVerified, string(profileData), u.UpdatedAt, u.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domerrors.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(id string) error {
	const q = `DELETE FROM users WHERE id = ?`

	result, err := r.db.Exec(q, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domerrors.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) GetAll(limit, offset int) ([]*domuser.User, int, error) {
	// Get total count
	var total int
	const countQ = `SELECT COUNT(*) FROM users`
	if err := r.db.QueryRow(countQ).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get users with pagination
	const q = `SELECT id, email, password_hash, first_name, last_name, role, is_active, email_verified, 
	           profile_data, created_at, updated_at, last_login_at FROM users 
	           ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(q, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users, err := r.scanUsers(rows)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) GetByRole(role domuser.Role, limit, offset int) ([]*domuser.User, int, error) {
	// Get total count for this role
	var total int
	const countQ = `SELECT COUNT(*) FROM users WHERE role = ?`
	if err := r.db.QueryRow(countQ, string(role)).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get users with pagination
	const q = `SELECT id, email, password_hash, first_name, last_name, role, is_active, email_verified, 
	           profile_data, created_at, updated_at, last_login_at FROM users 
	           WHERE role = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(q, string(role), limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users, err := r.scanUsers(rows)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) GetActiveUsers(limit, offset int) ([]*domuser.User, int, error) {
	// Get total count for active users
	var total int
	const countQ = `SELECT COUNT(*) FROM users WHERE is_active = 1`
	if err := r.db.QueryRow(countQ).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get users with pagination
	const q = `SELECT id, email, password_hash, first_name, last_name, role, is_active, email_verified, 
	           profile_data, created_at, updated_at, last_login_at FROM users 
	           WHERE is_active = 1 ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(q, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users, err := r.scanUsers(rows)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	const q = `SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`
	var exists bool
	err := r.db.QueryRow(q, email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) ExistsByID(id string) (bool, error) {
	const q = `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`
	var exists bool
	err := r.db.QueryRow(q, id).Scan(&exists)
	return exists, err
}

func (r *UserRepository) UpdateLastLogin(id string) error {
	const q = `UPDATE users SET last_login_at = ? WHERE id = ?`

	result, err := r.db.Exec(q, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domerrors.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) ToggleActiveStatus(id string, isActive bool) error {
	const q = `UPDATE users SET is_active = ?, updated_at = ? WHERE id = ?`

	result, err := r.db.Exec(q, isActive, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domerrors.ErrUserNotFound
	}

	return nil
}

// Helper method to scan multiple users from rows
func (r *UserRepository) scanUsers(rows *sql.Rows) ([]*domuser.User, error) {
	var users []*domuser.User

	for rows.Next() {
		u := &domuser.User{}
		var role string
		var profileData sql.NullString
		var lastLoginAt sql.NullTime

		err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
			&role, &u.IsActive, &u.EmailVerified, &profileData, &u.CreatedAt, &u.UpdatedAt, &lastLoginAt)
		if err != nil {
			return nil, err
		}

		u.Role = domuser.Role(role)
		if lastLoginAt.Valid {
			u.LastLoginAt = &lastLoginAt.Time
		}

		if profileData.Valid {
			if err := json.Unmarshal([]byte(profileData.String), &u.Profile); err != nil {
				u.Profile = domuser.Profile{}
			}
		}

		users = append(users, u)
	}

	return users, rows.Err()
}
