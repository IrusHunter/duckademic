package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type Course struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	ManagerID   *uuid.UUID `db:"manager_id" json:"manager_id"`
	Slug        string     `db:"slug" json:"slug"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

func (c Course) String() string {
	parts := make([]string, 0, 10)

	parts = append(parts, fmt.Sprintf("id: %s", c.ID))

	if c.ManagerID != nil {
		parts = append(parts, fmt.Sprintf("manager_id: %s", c.ManagerID))
	}

	parts = append(parts, fmt.Sprintf("slug: %s", c.Slug))
	parts = append(parts, fmt.Sprintf("name: %s", c.Name))

	if c.Description != "" {
		parts = append(parts, fmt.Sprintf("description: %s", c.Description))
	}

	if !c.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", c.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", c.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("Course{%s}", strings.Join(parts, ", "))
}

func (Course) TableName() string {
	return "courses"
}
func (Course) EntityName() string {
	return "course"
}
