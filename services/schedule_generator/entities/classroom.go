package entities

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Classroom struct {
	ID       uuid.UUID `json:"id"`
	Number   string    `json:"number"`
	Capacity int       `json:"capacity"`
}

func (c Classroom) String() string {
	parts := make([]string, 0, 10)

	if c.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", c.ID))
	}

	parts = append(parts, fmt.Sprintf("number: %s", c.Number))
	parts = append(parts, fmt.Sprintf("capacity: %d", c.Capacity))

	return fmt.Sprintf("Classroom{%s}", strings.Join(parts, ", "))
}

func (c *Classroom) ValidateNumber() error {
	if len(c.Number) == 0 {
		return fmt.Errorf("number required")
	}
	return nil
}

func (c *Classroom) ValidateCapacity() error {
	if c.Capacity <= 0 {
		return fmt.Errorf("capacity must be greater than 0 (%d <= 0)", c.Capacity)
	}
	return nil
}
