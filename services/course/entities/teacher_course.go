package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type TeacherCourse struct {
	ID        uuid.UUID `db:"id" json:"id"`
	CourseID  uuid.UUID `db:"course_id" json:"course_id"`
	TeacherID uuid.UUID `db:"teacher_id" json:"teacher_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (tc TeacherCourse) String() string {
	parts := make([]string, 0, 10)

	if tc.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", tc.ID))
	}
	if tc.CourseID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("course_id: %s", tc.CourseID))
	}
	if tc.TeacherID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("teacher_id: %s", tc.TeacherID))
	}

	if !tc.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", tc.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", tc.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("TeacherCourse{%s}", strings.Join(parts, ", "))
}

func (TeacherCourse) TableName() string {
	return "teacher_courses"
}
func (TeacherCourse) EntityName() string {
	return "teacher_course"
}
