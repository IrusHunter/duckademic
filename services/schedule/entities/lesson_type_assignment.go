package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type LessonTypeAssignment struct {
	ID            uuid.UUID `db:"id" json:"id"`
	LessonTypeID  uuid.UUID `db:"lesson_type_id" json:"lesson_type_id"`
	DisciplineID  uuid.UUID `db:"discipline_id" json:"discipline_id"`
	RequiredHours int       `db:"required_hours" json:"required_hours"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

func (lta LessonTypeAssignment) String() string {
	parts := make([]string, 0, 10)

	if lta.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", lta.ID))
	}

	parts = append(parts, fmt.Sprintf("lesson_type_id: %s", lta.LessonTypeID))
	parts = append(parts, fmt.Sprintf("discipline_id: %s", lta.DisciplineID))
	parts = append(parts, fmt.Sprintf("required_hours: %d", lta.RequiredHours))

	if !lta.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", lta.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", lta.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("LessonTypeAssignment{%s}", strings.Join(parts, ", "))
}

func (lta *LessonTypeAssignment) ValidateRequiredHours() error {
	if lta.RequiredHours <= 0 {
		return fmt.Errorf("required hours should be positive (%d <= 0)", lta.RequiredHours)
	}
	return nil
}

func (LessonTypeAssignment) TableName() string {
	return "lesson_type_assignments"
}
