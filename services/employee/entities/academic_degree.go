package entities

import (
	"fmt"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

// AcademicDegree represents a teacher's academic degree.
type AcademicDegree struct {
	ID        uuid.UUID `db:"id" json:"id"`                 // Unique identifier.
	Slug      string    `db:"slug" json:"slug"`             // Unique slug used internally.
	Title     string    `db:"title" json:"title"`           // Human-readable name of the degree.
	CreatedAt time.Time `db:"created_at" json:"created_at"` // Record creation timestamp.
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"` // Record last update timestamp.
}

// String returns a human-readable representation of the AcademicDegree.
// Includes title and optional ID, slug, created and updated timestamps.
func (ad *AcademicDegree) String() string {
	var createdAtStr, updatedAtStr, idStr, slugStr string
	if uuid.Nil != ad.ID {
		idStr = fmt.Sprintf("id: %s", ad.ID)
	}
	if ad.Slug != "" {
		slugStr = fmt.Sprintf(", slug: %s,", ad.Slug)
	}
	if !ad.CreatedAt.IsZero() {
		createdAtStr = fmt.Sprintf("created_at: %s", ad.CreatedAt.Format(db.TimeFormat))
		createdAtStr = fmt.Sprintf("created_at: %s", ad.CreatedAt.Format(db.TimeFormat))
	}
	return fmt.Sprintf("AcademicDegree{%s%s title: %s%s%s}",
		idStr, slugStr, ad.Title, createdAtStr, updatedAtStr,
	)
}

// ValidateTitle checks that Title is not empty.
func (ad *AcademicDegree) ValidateTitle() error {
	if len(ad.Title) == 0 {
		return fmt.Errorf("title required")
	}
	return nil
}
