package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	RecipientID   uuid.UUID  `db:"recipient_id" json:"recipient_id"`
	Type          string     `db:"type" json:"type"`
	Title         string     `db:"title" json:"title"`
	TaskID        uuid.UUID  `db:"task_id" json:"task_id"`
	CourseID      uuid.UUID  `db:"course_id" json:"course_id"`
	TaskStudentID *uuid.UUID `db:"task_student_id" json:"task_student_id,omitempty"`
	IsRead        bool       `db:"is_read" json:"is_read"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
}

func (n Notification) String() string {
	return fmt.Sprintf("Notification{id: %s, recipient: %s, type: %s}", n.ID, n.RecipientID, n.Type)
}

func (Notification) TableName() string  { return "notifications" }
func (Notification) EntityName() string { return "notification" }
