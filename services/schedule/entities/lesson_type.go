package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type LessonType struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Slug       string    `db:"slug" json:"slug"`
	Name       string    `db:"name" json:"name"`
	HoursValue int       `db:"hours_value" json:"hours_value"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

func (lt LessonType) String() string {
	parts := make([]string, 0, 10)

	if lt.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", lt.ID))
	}
	if lt.Slug != "" {
		parts = append(parts, fmt.Sprintf("slug: %s", lt.Slug))
	}

	parts = append(parts, fmt.Sprintf("name: %s", lt.Name))

	parts = append(parts, fmt.Sprintf("value: %d", lt.HoursValue))

	if !lt.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", lt.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", lt.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("LessonType{%s}", strings.Join(parts, ", "))
}

func (lt LessonType) TableName() string {
	return "lesson_types"
}
