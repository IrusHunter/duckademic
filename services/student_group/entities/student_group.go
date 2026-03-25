package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type StudentGroup struct {
	ID            uuid.UUID `db:"id" json:"id"`
	Slug          string    `db:"slug" json:"slug"`
	Name          string    `db:"name" json:"name"`
	GroupCohortID uuid.UUID `db:"group_cohort_id" json:"group_cohort_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

func (sg StudentGroup) String() string {
	parts := make([]string, 0, 10)

	if sg.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", sg.ID))
	}
	if sg.Slug != "" {
		parts = append(parts, fmt.Sprintf("slug: %s", sg.Slug))
	}

	parts = append(parts, fmt.Sprintf("name: %s", sg.Name))
	parts = append(parts, fmt.Sprintf("group_cohort_id: %s", sg.GroupCohortID))

	if !sg.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", sg.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", sg.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("StudentGroup{%s}", strings.Join(parts, ", "))
}

func (sg *StudentGroup) ValidateName() error {
	if len(sg.Name) == 0 {
		return fmt.Errorf("name required")
	}
	return nil
}

func (StudentGroup) TableName() string {
	return "student_groups"
}
