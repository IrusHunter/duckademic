package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type Discipline struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Slug      string    `db:"slug" json:"slug"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (d Discipline) String() string {
	parts := make([]string, 0, 6)

	if d.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", d.ID))
	}
	if d.Slug != "" {
		parts = append(parts, fmt.Sprintf("slug: %s", d.Slug))
	}

	parts = append(parts, fmt.Sprintf("name: %s", d.Name))

	if !d.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", d.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", d.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("Discipline{%s}", strings.Join(parts, ", "))
}

func (d Discipline) TableName() string {
	return "disciplines"
}
