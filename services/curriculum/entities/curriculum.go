package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type Curriculum struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	Slug          string     `db:"slug" json:"slug"`
	Name          string     `db:"name" json:"name"`
	DurationYears int        `db:"duration_years" json:"duration_years"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
	EffectiveFrom time.Time  `db:"effective_from" json:"effective_from"`
	EffectiveTo   *time.Time `db:"effective_to" json:"effective_to,omitempty"`
}

func (c Curriculum) String() string {
	parts := make([]string, 0, 10)

	if c.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", c.ID))
	}
	if c.Slug != "" {
		parts = append(parts, fmt.Sprintf("slug: %s", c.Slug))
	}

	parts = append(parts, fmt.Sprintf("name: %s", c.Name))
	parts = append(parts, fmt.Sprintf("duration_years: %d", c.DurationYears))

	if !c.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", c.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", c.UpdatedAt.Format(db.TimeFormat)))
	}

	parts = append(parts, fmt.Sprintf("effective_from: %s", c.EffectiveFrom.Format(db.TimeFormat)))

	if c.EffectiveTo != nil {
		parts = append(parts, fmt.Sprintf("effective_to: %s", c.EffectiveTo.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("Curriculum{%s}", strings.Join(parts, ", "))
}

func (c *Curriculum) ValidateName() error {
	if len(c.Name) == 0 {
		return fmt.Errorf("name required")
	}
	return nil
}

func (c *Curriculum) ValidateDurationYears() error {
	if c.DurationYears <= 0 {
		return fmt.Errorf("should last longer then 0 years (%d <= 0)", c.DurationYears)
	}
	return nil
}

func (c *Curriculum) ValidateEffectiveFrom() error {
	if c.EffectiveFrom.IsZero() {
		return fmt.Errorf("effective from timestamp required")
	}
	return nil
}

func (c Curriculum) TableName() string {
	return "curriculums"
}
