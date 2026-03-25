package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type GroupCohort struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Slug       string    `db:"slug" json:"slug"`
	Name       string    `db:"name" json:"name"`
	SemesterID uuid.UUID `db:"semester_id" json:"semester_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

func (g GroupCohort) String() string {
	parts := make([]string, 0, 10)

	if g.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", g.ID))
	}
	if g.Slug != "" {
		parts = append(parts, fmt.Sprintf("slug: %s", g.Slug))
	}

	parts = append(parts, fmt.Sprintf("name: %s", g.Name))
	parts = append(parts, fmt.Sprintf("semester_id: %s", g.SemesterID))

	if !g.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", g.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", g.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("GroupCohort{%s}", strings.Join(parts, ", "))
}

func (gc *GroupCohort) ValidateName() error {
	if len(gc.Name) == 0 {
		return fmt.Errorf("name required")
	}

	return nil
}

func (GroupCohort) TableName() string {
	return "group_cohorts"
}
