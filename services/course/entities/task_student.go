package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type TaskStudent struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	TaskID         uuid.UUID  `db:"task_id" json:"task_id"`
	StudentID      uuid.UUID  `db:"student_id" json:"student_id"`
	Mark           *float64   `db:"mark" json:"mark,omitempty"`
	SubmissionTime *time.Time `db:"submission_time" json:"submission_time,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
}

func (ts TaskStudent) String() string {
	parts := make([]string, 0, 12)

	if ts.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", ts.ID))
	}
	if ts.TaskID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("task_id: %s", ts.TaskID))
	}
	if ts.StudentID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("student_id: %s", ts.StudentID))
	}
	if ts.Mark != nil {
		parts = append(parts, fmt.Sprintf("mark: %f", *ts.Mark))
	}
	if ts.SubmissionTime != nil && !ts.SubmissionTime.IsZero() {
		parts = append(parts, fmt.Sprintf("submission_time: %s", ts.SubmissionTime.Format(db.TimeFormat)))
	}
	if !ts.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", ts.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", ts.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("TaskStudent{%s}", strings.Join(parts, ", "))
}

func (TaskStudent) TableName() string {
	return "task_students"
}
func (TaskStudent) EntityName() string {
	return "task student"
}
