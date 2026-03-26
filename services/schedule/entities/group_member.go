package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type GroupMember struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	StudentID    uuid.UUID  `db:"student_id" json:"studentId"`
	StudentGroup *uuid.UUID `db:"student_group_id" json:"student_group_id,omitempty"`
	CreatedAt    time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updatedAt"`
}

func (g GroupMember) String() string {
	parts := make([]string, 0, 10)

	if g.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", g.ID))
	}

	parts = append(parts, fmt.Sprintf("student_id: %s", g.StudentID))

	if g.StudentGroup != nil {
		parts = append(parts, fmt.Sprintf("student_group_id: %s", *g.StudentGroup))
	}

	if !g.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", g.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", g.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("GroupMembers{%s}", strings.Join(parts, ", "))
}

func (GroupMember) TableName() string {
	return "group_members"
}
