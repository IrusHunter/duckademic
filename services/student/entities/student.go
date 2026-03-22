package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type Student struct {
	ID          uuid.UUID  `db:"id" json:"id"`                               // Unique identifier.
	Slug        string     `db:"slug" json:"slug"`                           // Unique slug used internally.
	FirstName   string     `db:"first_name" json:"first_name"`               // Employees first name.
	LastName    string     `db:"last_name" json:"last_name"`                 // Employees last name.
	MiddleName  *string    `db:"middle_name" json:"middle_name,omitempty"`   // Employees middle name.
	PhoneNumber *string    `db:"phone_number" json:"phone_number,omitempty"` // Contact phone number.
	Email       string     `db:"email" json:"email"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`           // Record creation timestamp.
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`           // Record last update timestamp.
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"` // Record deleted timestamp.
}

func (s Student) String() string {
	parts := make([]string, 0, 10)

	if s.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", s.ID))
	}
	if s.Slug != "" {
		parts = append(parts, fmt.Sprintf("slug: %s", s.Slug))
	}

	parts = append(parts, fmt.Sprintf("first_name: %s", s.FirstName))
	parts = append(parts, fmt.Sprintf("last_name: %s", s.LastName))

	if s.MiddleName != nil {
		parts = append(parts, fmt.Sprintf("middle_name: %s", *s.MiddleName))
	}
	if s.PhoneNumber != nil {
		parts = append(parts, fmt.Sprintf("phone_number: %s", *s.PhoneNumber))
	}

	parts = append(parts, fmt.Sprintf("email: %s", s.Email))

	if !s.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", s.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", s.UpdatedAt.Format(db.TimeFormat)))
	}
	if s.DeletedAt != nil {
		parts = append(parts, fmt.Sprintf("deleted_at: %s", s.DeletedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("Student{%s}", strings.Join(parts, ", "))
}

// ValidateFirstName checks that FirstName is not empty.
func (s *Student) ValidateFirstName() error {
	if len(s.FirstName) == 0 {
		return fmt.Errorf("first name required")
	}
	return nil
}

// ValidateLastName checks that LastName is not empty.
func (s *Student) ValidateLastName() error {
	if len(s.LastName) == 0 {
		return fmt.Errorf("last name required")
	}
	return nil
}

// ValidateEmail checks that Email is not empty.
func (s *Student) ValidateEmail() error {
	if len(s.Email) == 0 {
		return fmt.Errorf("email required")
	}
	return nil
}

// GetFullName returns full name (Last First Middle).
func (s *Student) GetFullName() string {
	parts := make([]string, 0, 3)

	parts = append(parts, s.LastName)
	parts = append(parts, s.FirstName)

	if s.MiddleName != nil {
		parts = append(parts, *s.MiddleName)
	}

	return strings.Join(parts, " ")
}

// GetShortFullName returns short name (Last F.M.).
func (s *Student) GetShortFullName() string {
	parts := make([]string, 0, 3)

	parts = append(parts, fmt.Sprintf("%s ", s.LastName))
	parts = append(parts, fmt.Sprintf("%c.", []rune(s.FirstName)[0]))

	if s.MiddleName != nil {
		parts = append(parts, fmt.Sprintf("%c.", []rune(*s.MiddleName)[0]))
	}

	return strings.Join(parts, "")
}

func (Student) TableName() string {
	return "students"
}
