package user

import "time"

type Role string

const (
	RoleUser      Role = "user"
	RoleOperator  Role = "operator"
	RoleAdmin     Role = "admin"
	RoleTreasurer Role = "treasurer"
)

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Role          Role      `json:"role"`
	IsActive      bool      `json:"isActive"`
	EmailVerified bool      `json:"emailVerified"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
