package entities

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type LessonTypeAssignment struct {
	ID            uuid.UUID `json:"id"`
	LessonTypeID  uuid.UUID `json:"lesson_type_id"`
	DisciplineID  uuid.UUID `json:"discipline_id"`
	RequiredHours int       `json:"required_hours"`
}

func (lta LessonTypeAssignment) String() string {
	parts := make([]string, 0, 10)

	parts = append(parts, fmt.Sprintf("lesson_type_id: %s", lta.LessonTypeID))
	parts = append(parts, fmt.Sprintf("discipline_id: %s", lta.DisciplineID))
	parts = append(parts, fmt.Sprintf("required_hours: %d", lta.RequiredHours))

	return fmt.Sprintf("LessonTypeAssignment{%s}", strings.Join(parts, ", "))
}

func (lta *LessonTypeAssignment) ValidateRequiredHours() error {
	if lta.RequiredHours <= 0 {
		return fmt.Errorf("required hours should be positive (%d <= 0)", lta.RequiredHours)
	}
	return nil
}
