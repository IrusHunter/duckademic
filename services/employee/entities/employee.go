package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

// Employee represented an university employee.
type Employee struct {
	ID          uuid.UUID  `db:"id" json:"id"`                               // Unique identifier.
	Slug        string     `db:"slug" json:"slug"`                           // Unique slug used internally.
	FirstName   string     `db:"first_name" json:"first_name"`               // Employees first name.
	LastName    string     `db:"last_name" json:"last_name"`                 // Employees last name.
	MiddleName  *string    `db:"middle_name" json:"middle_name,omitempty"`   // Employees middle name.
	PhoneNumber *string    `db:"phone_number" json:"phone_number,omitempty"` // Contact phone number.
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`               // Record creation timestamp.
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`               // Record last update timestamp.
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`     // Record deleted timestamp.
}

// String returns a human-readable representation of the Employee.
// Includes first and last names and optional ID, slug, middle name, phone number, created,
// updated and deleted timestamps.
func (e Employee) String() string {
	parts := make([]string, 0, 10)
	if e.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", e.ID))
	}
	if e.Slug != "" {
		parts = append(parts, fmt.Sprintf("slug: %s", e.Slug))
	}
	parts = append(parts, fmt.Sprintf("first_name: %s", e.FirstName))
	parts = append(parts, fmt.Sprintf("last_name: %s", e.LastName))
	if e.MiddleName != nil {
		parts = append(parts, fmt.Sprintf("middle_name: %s", *e.MiddleName))
	}
	if e.PhoneNumber != nil {
		parts = append(parts, fmt.Sprintf("phone_number: %s", *e.PhoneNumber))
	}
	if !e.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", e.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", e.UpdatedAt.Format(db.TimeFormat)))
	}
	if e.DeletedAt != nil {
		parts = append(parts, fmt.Sprintf("deleted_at: %s", e.DeletedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("Employee{%s}", strings.Join(parts, ", "))
}

// ValidateFirstName checks that FirstName is not empty.
func (e *Employee) ValidateFirstName() error {
	if len(e.FirstName) == 0 {
		return fmt.Errorf("first name required")
	}
	return nil
}

// ValidateLastName checks that LastName is not empty.
func (e *Employee) ValidateLastName() error {
	if len(e.LastName) == 0 {
		return fmt.Errorf("last name required")
	}
	return nil
}

// GetFullName returns the employee's full name. It includes the first name, last name,
// and, optionally, the middle name.
func (e *Employee) GetFullName() string {
	parts := make([]string, 0, 3)
	parts = append(parts, e.LastName)
	parts = append(parts, e.FirstName)
	if e.MiddleName != nil {
		parts = append(parts, *e.MiddleName)
	}
	return strings.Join(parts, " ")
}

func (e *Employee) GetShortFullName() string {
	parts := make([]string, 0, 3)
	parts = append(parts, fmt.Sprintf("%s ", e.LastName))
	parts = append(parts, fmt.Sprintf("%s.", string(e.FirstName[0])))
	if e.MiddleName != nil {
		parts = append(parts, fmt.Sprintf("%s.", string((*e.MiddleName)[0])))
	}
	return strings.Join(parts, "")
}

func (Employee) TableName() string {
	return "employees"
}
