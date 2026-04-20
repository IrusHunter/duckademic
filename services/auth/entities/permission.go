package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type Permission struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (p Permission) String() string {
	parts := make([]string, 0, 4)

	if p.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", p.ID))
	}

	parts = append(parts, fmt.Sprintf("name: %s", p.Name))

	if !p.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", p.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", p.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("Permission{%s}", strings.Join(parts, ", "))
}
func (Permission) TableName() string {
	return "permissions"
}
func (Permission) EntityName() string {
	return "permission"
}
