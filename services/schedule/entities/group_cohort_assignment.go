package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type GroupCohortAssignment struct {
	ID            uuid.UUID `db:"id" json:"id"`
	GroupCohortID uuid.UUID `db:"group_cohort_id" json:"group_cohort_id"`
	DisciplineID  uuid.UUID `db:"discipline_id" json:"discipline_id"`
	LessonTypeID  uuid.UUID `db:"lesson_type_id" json:"lesson_type_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

func (g GroupCohortAssignment) String() string {
	parts := make([]string, 0, 6)

	if g.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", g.ID))
	}

	parts = append(parts, fmt.Sprintf("group_cohort_id: %s", g.GroupCohortID))
	parts = append(parts, fmt.Sprintf("discipline_id: %s", g.DisciplineID))
	parts = append(parts, fmt.Sprintf("lesson_type_id: %s", g.LessonTypeID))

	if !g.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", g.CreatedAt.Format(time.RFC3339)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", g.UpdatedAt.Format(time.RFC3339)))
	}

	return fmt.Sprintf("GroupCohortAssignment{%s}", strings.Join(parts, ", "))
}

func (GroupCohortAssignment) TableName() string {
	return "group_cohort_assignments"
}
