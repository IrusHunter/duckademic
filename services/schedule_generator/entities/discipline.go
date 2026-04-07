package entities

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Discipline struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (d Discipline) String() string {
	parts := make([]string, 0, 6)

	parts = append(parts, fmt.Sprintf("id: %s", d.ID))
	parts = append(parts, fmt.Sprintf("name: %s", d.Name))

	return fmt.Sprintf("Discipline{%s}", strings.Join(parts, ", "))
}

func (d *Discipline) ValidateName() error {
	if len(d.Name) == 0 {
		return fmt.Errorf("name required")
	}
	return nil
}
