package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type StudyLoad struct {
	ID               uuid.UUID `json:"id" db:"id"`
	TeacherID        uuid.UUID `json:"teacher_id" db:"teacher_id"`
	TeacherName      string    `json:"teacher_name" db:"teacher_name"`
	StudentGroupID   uuid.UUID `json:"student_group_id" db:"student_group_id"`
	StudentGroupName string    `json:"student_group_name" db:"student_group_name"`
	DisciplineID     uuid.UUID `json:"discipline_id" db:"discipline_id"`
	DisciplineName   string    `json:"discipline_name" db:"discipline_name"`
	LessonTypeID     uuid.UUID `json:"lesson_type_id" db:"lesson_type_id"`
	LessonTypeName   string    `json:"lesson_type_name" db:"lesson_type_name"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

func (s StudyLoad) String() string {
	parts := make([]string, 0, 16)

	if s.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", s.ID))
	}

	parts = append(parts, fmt.Sprintf("teacher_id: %s", s.TeacherID))
	parts = append(parts, fmt.Sprintf("teacher_name: %s", s.TeacherName))
	parts = append(parts, fmt.Sprintf("student_group_id: %s", s.StudentGroupID))
	parts = append(parts, fmt.Sprintf("student_group_name: %s", s.StudentGroupName))
	parts = append(parts, fmt.Sprintf("discipline_id: %s", s.DisciplineID))
	parts = append(parts, fmt.Sprintf("discipline_name: %s", s.DisciplineName))
	parts = append(parts, fmt.Sprintf("lesson_type_id: %s", s.LessonTypeID))
	parts = append(parts, fmt.Sprintf("lesson_type_name: %s", s.LessonTypeName))

	if !s.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", s.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", s.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("StudyLoad{%s}", strings.Join(parts, ", "))
}

func (StudyLoad) TableName() string {
	return "study_loads"
}
func (StudyLoad) EntityName() string {
	return "study load"
}
