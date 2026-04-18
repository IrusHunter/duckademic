package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type StudyLoad struct {
	ID             uuid.UUID `json:"id" db:"id"`
	TeacherID      uuid.UUID `json:"teacher_id" db:"teacher_id"`
	StudentGroupID uuid.UUID `json:"student_group_id" db:"student_group_id"`
	DisciplineID   uuid.UUID `json:"discipline_id" db:"discipline_id"`
	LessonTypeID   uuid.UUID `json:"lesson_type_id" db:"lesson_type_id"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

func (s StudyLoad) String() string {
	parts := make([]string, 0, 10)

	if s.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", s.ID))
	}

	parts = append(parts, fmt.Sprintf("teacher_id: %s", s.TeacherID))
	parts = append(parts, fmt.Sprintf("student_group_id: %s", s.StudentGroupID))
	parts = append(parts, fmt.Sprintf("discipline_id: %s", s.DisciplineID))
	parts = append(parts, fmt.Sprintf("lesson_type_id: %s", s.LessonTypeID))

	if !s.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", s.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", s.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("StudyLoad{%s}", strings.Join(parts, ", "))
}

func (StudyLoad) TableName() string {
	return "study_loads"
}

type CompactStudyLoad struct {
	TeacherID        uuid.UUID
	TeacherName      string
	StudentGroupID   uuid.UUID
	StudentGroupName string
	DisciplineID     uuid.UUID
	DisciplineName   string
	LessonTypeID     uuid.UUID
	LessonTypeName   string
}
