package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type Teacher struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	Slug           string     `db:"slug" json:"slug"`
	Name           string     `db:"name" json:"name"`
	AcademicRankID uuid.UUID  `db:"academic_rank_id" json:"academic_rank_id"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`           // Record creation timestamp.
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`           // Record last update timestamp.
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at,omitempty"` // Record deleted timestamp.
}

func (t Teacher) String() string {
	parts := make([]string, 0, 10)
	if t.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", t.ID))
	}
	if t.Slug != "" {
		parts = append(parts, fmt.Sprintf("slug: %s", t.Slug))
	}
	parts = append(parts, fmt.Sprintf("name: %s", t.Name))
	parts = append(parts, fmt.Sprintf("academic rank id: %s", t.AcademicRankID))
	if !t.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", t.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", t.UpdatedAt.Format(db.TimeFormat)))
	}
	if t.DeletedAt != nil {
		parts = append(parts, fmt.Sprintf("deleted_at: %s", t.DeletedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("Teacher{%s}", strings.Join(parts, ", "))
}

func (t *Teacher) ValidateName() error {
	if t.Name == "" {
		return fmt.Errorf("name required")
	}

	return nil
}

func (Teacher) TableName() string {
	return "teachers"
}
