package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID `db:"id" json:"id"`
	CourseID    uuid.UUID `db:"course_id" json:"course_id"`
	Slug        string    `db:"slug" json:"slug"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	MaxMark     float64   `db:"max_mark" json:"max_mark"`
	Deadline    time.Time `db:"deadline" json:"deadline"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func (t Task) String() string {
	parts := make([]string, 0, 12)

	if t.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", t.ID))
	}

	parts = append(parts, fmt.Sprintf("course_id: %s", t.CourseID))
	parts = append(parts, fmt.Sprintf("slug: %s", t.Slug))
	parts = append(parts, fmt.Sprintf("title: %s", t.Title))
	parts = append(parts, fmt.Sprintf("description: %s", t.Description))
	parts = append(parts, fmt.Sprintf("max_mark: %f", t.MaxMark))
	parts = append(parts, fmt.Sprintf("deadline: %s", t.Deadline.Format(db.TimeFormat)))

	if !t.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", t.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", t.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("Task{%s}", strings.Join(parts, ", "))
}

func (Task) TableName() string {
	return "tasks"
}
func (Task) EntityName() string {
	return "task"
}
