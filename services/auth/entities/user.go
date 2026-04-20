package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `db:"id" json:"id"`
	Login             string    `db:"login" json:"login"`
	HashedPassword    string    `db:"password" json:"-"`
	IsDefaultPassword bool      `db:"is_default_password" json:"is_default_password"`
	RoleID            uuid.UUID `db:"role_id" json:"role_id"`
	LastLogin         time.Time `db:"last_login" json:"last_login"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`

	Password *string `db:"-" json:"password,omitempty"`
	RoleName *string `db:"-" json:"role,omitempty"`
}

func (u User) String() string {
	parts := make([]string, 0, 6)

	if u.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", u.ID))
	}

	parts = append(parts, fmt.Sprintf("login: %s", u.Login))
	parts = append(parts, fmt.Sprintf("is_default_password: %t", u.IsDefaultPassword))

	if u.RoleID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("role_id: %s", u.RoleID))
	}

	if !u.LastLogin.IsZero() {
		parts = append(parts, fmt.Sprintf("last_login: %s", u.LastLogin.Format(db.TimeFormat)))
	}

	if !u.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", u.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", u.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("User{%s}", strings.Join(parts, ", "))
}
func (User) TableName() string {
	return "users"
}
func (User) EntityName() string {
	return "user"
}

func (u *User) GetPassword() string {
	if u.Password == nil {
		return ""
	}
	return *u.Password
}
