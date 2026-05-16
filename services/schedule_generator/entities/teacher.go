package entities

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type Teacher struct {
	ID              uuid.UUID         `json:"id"`
	Name            string            `json:"name"`
	Priority        int               `json:"priority"`
	SlotsPriorities map[int][]float32 `json:"slots_priorities"`
	UnavailableDays []time.Time       `json:"unavailable_days"`
}

func (t Teacher) String() string {
	parts := make([]string, 0, 6)

	parts = append(parts, fmt.Sprintf("id: %s", t.ID.String()))
	parts = append(parts, fmt.Sprintf("name: %s", t.Name))
	parts = append(parts, fmt.Sprintf("priority: %d", t.Priority))

	slotParts := make([]string, 0, len(t.SlotsPriorities))
	for slot, priorities := range t.SlotsPriorities {
		priorityStrs := make([]string, len(priorities))
		for i, p := range priorities {
			priorityStrs[i] = fmt.Sprintf("%.2f", p)
		}

		slotParts = append(
			slotParts,
			fmt.Sprintf("%d:[%s]", slot, strings.Join(priorityStrs, ", ")),
		)
	}
	parts = append(parts,
		fmt.Sprintf("slot_priorities: {%s}", strings.Join(slotParts, ", ")),
	)

	dayParts := make([]string, len(t.UnavailableDays))
	for i, day := range t.UnavailableDays {
		dayParts[i] = day.Format(db.TimeFormat)
	}
	parts = append(parts,
		fmt.Sprintf("unavailable_days: [%s]", strings.Join(dayParts, ", ")),
	)

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
func (t *Teacher) ValidateUnavailableDays() error {
	for i, day := range t.UnavailableDays {
		if day.IsZero() {
			return fmt.Errorf("zero day at %d", i)
		}
	}

	return nil
}
