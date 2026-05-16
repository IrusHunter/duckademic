package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type TeacherUnavailableDay struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TeacherID uuid.UUID `json:"teacher_id" db:"teacher_id"`
	Day       time.Time `json:"day" db:"day"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (t TeacherUnavailableDay) String() string {
	parts := make([]string, 0, 6)

	if t.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", t.ID))
	}

	parts = append(parts, fmt.Sprintf("teacher_id: %s", t.TeacherID))
	parts = append(parts, fmt.Sprintf("day: %s", t.Day.Format(time.DateOnly)))

	if !t.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", t.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", t.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("TeacherUnavailableDays{%s}", strings.Join(parts, ", "))
}

func (t *TeacherUnavailableDay) ValidateDay() error {
	if t.Day.IsZero() {
		return fmt.Errorf("day required")
	}

	return nil
}

func (TeacherUnavailableDay) TableName() string {
	return "teacher_unavailable_days"
}

func (TeacherUnavailableDay) EntityName() string {
	return "teacher unavailable day"
}
