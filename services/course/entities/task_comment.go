package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TaskComment struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	TaskID     uuid.UUID  `db:"task_id" json:"task_id"`
	AuthorID   uuid.UUID  `db:"author_id" json:"author_id"`
	AuthorName string     `db:"author_name" json:"author_name"`
	Body       string     `db:"body" json:"body"`
	IsPrivate  bool       `db:"is_private" json:"is_private"`
	StudentID  *uuid.UUID `db:"student_id" json:"student_id,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
}

func (c TaskComment) String() string {
	parts := make([]string, 0, 6)
	if c.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", c.ID))
	}
	if c.TaskID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("task_id: %s", c.TaskID))
	}
	if c.AuthorName != "" {
		parts = append(parts, fmt.Sprintf("author: %s", c.AuthorName))
	}
	if c.Body != "" {
		parts = append(parts, fmt.Sprintf("body: %s", strings.TrimSpace(c.Body)))
	}
	return fmt.Sprintf("TaskComment{%s}", strings.Join(parts, ", "))
}

func (TaskComment) TableName() string { return "task_comments" }
func (TaskComment) EntityName() string { return "task comment" }
