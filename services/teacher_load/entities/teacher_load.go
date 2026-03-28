package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type TeacherLoad struct {
	ID            uuid.UUID `db:"id" json:"id"`
	TeacherID     uuid.UUID `db:"teacher_id" json:"teacher_id"`
	DisciplineID  uuid.UUID `db:"discipline_id" json:"discipline_id"`
	LessonTypeID  uuid.UUID `db:"lesson_type_id" json:"lesson_type_id"`
	GroupCohortID uuid.UUID `db:"group_cohort_id" json:"group_cohort_id"`
	GroupCount    int       `db:"group_count" json:"group_count"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

func (t TeacherLoad) String() string {
	parts := make([]string, 0, 10)

	if t.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", t.ID))
	}

	parts = append(parts, fmt.Sprintf("teacher_id: %s", t.TeacherID))
	parts = append(parts, fmt.Sprintf("discipline_id: %s", t.DisciplineID))
	parts = append(parts, fmt.Sprintf("lesson_type_id: %s", t.LessonTypeID))
	parts = append(parts, fmt.Sprintf("group_cohort_id: %s", t.GroupCohortID))
	parts = append(parts, fmt.Sprintf("group_count: %d", t.GroupCount))

	if !t.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", t.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", t.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("TeacherLoad{%s}", strings.Join(parts, ", "))
}

func (t *TeacherLoad) ValidateGroupCount() error {
	if t.GroupCount <= 0 {
		return fmt.Errorf("group count should be positive (%d <= 0)", t.GroupCount)
	}
	return nil
}

func (TeacherLoad) TableName() string {
	return "teacher_loads"
}
