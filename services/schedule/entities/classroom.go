package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type Classroom struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Slug      string    `db:"slug" json:"slug"`
	Number    string    `db:"number" json:"number"`
	Capacity  int       `db:"capacity" json:"capacity"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (c Classroom) String() string {
	parts := make([]string, 0, 10)

	if c.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", c.ID))
	}
	if c.Slug != "" {
		parts = append(parts, fmt.Sprintf("slug: %s", c.Slug))
	}

	parts = append(parts, fmt.Sprintf("number: %s", c.Number))
	parts = append(parts, fmt.Sprintf("capacity: %d", c.Capacity))

	if !c.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", c.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", c.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("Classroom{%s}", strings.Join(parts, ", "))
}

func (c Classroom) TableName() string {
	return "classrooms"
}
