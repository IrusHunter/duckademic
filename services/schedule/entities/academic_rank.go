package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

// AcademicRank represents the teacher's academic rank.
type AcademicRank struct {
	ID        uuid.UUID `db:"id" json:"id"`                 // Unique identifier.
	Slug      string    `db:"slug" json:"slug"`             // Unique slug used internally.
	Title     string    `db:"title" json:"title"`           // Human-readable name of the rank.
	Priority  int       `db:"priority" json:"priority"`     // Determines the rank's priority: higher value = higher rank.
	CreatedAt time.Time `db:"created_at" json:"created_at"` // Record creation timestamp.
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"` // Record last update timestamp.
}

// String returns a human-readable representation of the AcademicRank.
// Includes title, priority and optional ID, slug, created and updated timestamps.
func (ar AcademicRank) String() string {
	parts := make([]string, 0, 6)

	if ar.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", ar.ID))
	}
	if ar.Slug != "" {
		parts = append(parts, fmt.Sprintf("slug: %s", ar.Slug))
	}
	parts = append(parts, fmt.Sprintf("title: %s", ar.Title))
	parts = append(parts, fmt.Sprintf("priority: %d", ar.Priority))

	if !ar.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", ar.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", ar.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("AcademicRank{%s}", strings.Join(parts, ", "))
}

// ValidateTitle checks that Title is not empty.
func (ar *AcademicRank) ValidateTitle() error {
	if len(ar.Title) == 0 {
		return fmt.Errorf("title required")
	}
	return nil
}

func (AcademicRank) TableName() string {
	return "academic_ranks"
}
