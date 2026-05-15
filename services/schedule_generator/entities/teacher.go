package entities

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Teacher struct {
	ID              uuid.UUID         `json:"id"`
	Name            string            `json:"name"`
	Priority        int               `json:"priority"`
	SlotsPriorities map[int][]float32 `json:"slots_priorities"`
}

func (t Teacher) String() string {
	parts := make([]string, 0, 3)

	parts = append(parts, fmt.Sprintf("id: %s", t.ID.String()))
	parts = append(parts, fmt.Sprintf("name: %s", t.Name))
	parts = append(parts, fmt.Sprintf("priority: %d", t.Priority))

	return fmt.Sprintf("Teacher{%s}", strings.Join(parts, ", "))
}

func (t *Teacher) ValidateName() error {
	if t.Name == "" {
		return errors.New("name must not be empty")
	}
	return nil
}
func (t *Teacher) ValidateSlotsPriorities() error {
	for day, priorities := range t.SlotsPriorities {
		for i, priority := range priorities {
			if priority < 0 {
				return fmt.Errorf("priority should be positive (invalid at [%d, %d])", day, i)
			}
		}
	}

	return nil
}
