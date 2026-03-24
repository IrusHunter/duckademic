package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type Semester struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Slug         string    `db:"slug" json:"slug"`
	CurriculumID uuid.UUID `db:"curriculum_id" json:"curriculum_id"`
	Number       int       `db:"number" json:"number"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

func (s Semester) String() string {
	parts := make([]string, 0, 10)

	if s.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", s.ID))
	}
	if s.Slug != "" {
		parts = append(parts, fmt.Sprintf("slug: %s", s.Slug))
	}

	parts = append(parts, fmt.Sprintf("curriculum_id: %s", s.CurriculumID))
	parts = append(parts, fmt.Sprintf("number: %d", s.Number))

	if !s.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", s.CreatedAt.Format(db.TimeFormat)))
	}
	if !s.UpdatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("updated_at: %s", s.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("Semester{%s}", strings.Join(parts, ", "))
}

func (s *Semester) ValidateNumber() error {
	if s.Number <= 0 {
		return fmt.Errorf("semester number should be positive (%d <= 0)", s.Number)
	}
	return nil
}

func (Semester) TableName() string {
	return "semesters"
}
