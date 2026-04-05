package entities

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type TeacherLoad struct {
	ID           uuid.UUID `json:"id"`
	TeacherID    uuid.UUID `json:"teacher_id"`
	DisciplineID uuid.UUID `json:"discipline_id"`
	LessonTypeID uuid.UUID `json:"lesson_type_id"`
	GroupCount   int       `json:"group_count"`
}

func (t TeacherLoad) String() string {
	parts := make([]string, 0, 10)

	if t.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", t.ID))
	}

	parts = append(parts, fmt.Sprintf("teacher_id: %s", t.TeacherID))
	parts = append(parts, fmt.Sprintf("discipline_id: %s", t.DisciplineID))
	parts = append(parts, fmt.Sprintf("lesson_type_id: %s", t.LessonTypeID))
	parts = append(parts, fmt.Sprintf("group_count: %d", t.GroupCount))

	return fmt.Sprintf("TeacherLoad{%s}", strings.Join(parts, ", "))
}
