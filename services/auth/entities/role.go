package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type Role struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (r Role) String() string {
	parts := make([]string, 0, 4)

	if r.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", r.ID))
	}

	parts = append(parts, fmt.Sprintf("name: %s", r.Name))

	if !r.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", r.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", r.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("Role{%s}", strings.Join(parts, ", "))
}

func (Role) TableName() string {
	return "roles"
}

func (Role) EntityName() string {
	return "role"
}
