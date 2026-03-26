package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type Student struct {
	ID        uuid.UUID  `db:"id" json:"id"`                           // Unique identifier.
	Slug      string     `db:"slug" json:"slug"`                       // Unique slug used internally.
	Name      string     `db:"name" json:"name"`                       // Employees first name.
	CreatedAt time.Time  `db:"created_at" json:"created_at"`           // Record creation timestamp.
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`           // Record last update timestamp.
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"` // Record deleted timestamp.
}

func (s Student) String() string {
	parts := make([]string, 0, 10)

	if s.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", s.ID))
	}
	if s.Slug != "" {
		parts = append(parts, fmt.Sprintf("slug: %s", s.Slug))
	}

	parts = append(parts, fmt.Sprintf("name: %s", s.Name))

	if !s.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", s.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", s.UpdatedAt.Format(db.TimeFormat)))
	}
	if s.DeletedAt != nil {
		parts = append(parts, fmt.Sprintf("deleted_at: %s", s.DeletedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("Student{%s}", strings.Join(parts, ", "))
}

func (Student) TableName() string {
	return "students"
}
