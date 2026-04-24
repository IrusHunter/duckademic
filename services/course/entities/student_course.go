package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type StudentCourse struct {
	ID        uuid.UUID `db:"id" json:"id"`
	CourseID  uuid.UUID `db:"course_id" json:"course_id"`
	StudentID uuid.UUID `db:"student_id" json:"student_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (sc StudentCourse) String() string {
	parts := make([]string, 0, 10)

	if sc.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", sc.ID))
	}

	parts = append(parts, fmt.Sprintf("course_id: %s", sc.CourseID))
	parts = append(parts, fmt.Sprintf("student_id: %s", sc.StudentID))

	if !sc.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", sc.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", sc.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("StudentCourse{%s}", strings.Join(parts, ", "))
}

func (StudentCourse) TableName() string {
	return "student_courses"
}

func (StudentCourse) EntityName() string {
	return "student_course"
}
