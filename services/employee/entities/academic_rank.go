package entities

import (
	"fmt"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

// AcademicRank represents the teacher's academic rank.
type AcademicRank struct {
	ID        uuid.UUID `db:"id" json:"id"`                 // Unique identifier.
	Slug      string    `db:"slug" json:"slug"`             // Unique slug used internally.
	Title     string    `db:"title" json:"title"`           // Human-readable name of the rank.
	CreatedAt time.Time `db:"created_at" json:"created_at"` // Record creation timestamp.
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"` // Record last update timestamp.
}

// String returns a human-readable representation of the AcademicRank.
// Includes title and optional ID, slug, created and updated timestamps.
func (ar AcademicRank) String() string {
	var createdAtStr, updatedAtStr, idStr, slugStr string
	if uuid.Nil != ar.ID {
		idStr = fmt.Sprintf("id: %s", ar.ID)
	}
	if ar.Slug != "" {
		slugStr = fmt.Sprintf(", slug: %s,", ar.Slug)
	}
	if !ar.CreatedAt.IsZero() {
		createdAtStr = fmt.Sprintf("created_at: %s", ar.CreatedAt.Format(db.TimeFormat))
		createdAtStr = fmt.Sprintf("created_at: %s", ar.CreatedAt.Format(db.TimeFormat))
	}
	return fmt.Sprintf("AcademicRank{%s%s title: %s%s%s}",
		idStr, slugStr, ar.Title, createdAtStr, updatedAtStr,
	)
}

// ValidateTitle checks that Title is not empty.
func (ar *AcademicRank) ValidateTitle() error {
	if len(ar.Title) == 0 {
		return fmt.Errorf("title required")
	}
	return nil
}
