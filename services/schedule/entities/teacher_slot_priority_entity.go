package entities

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type TeacherSlotPriorityValue string

const (
	TeacherSlotBlocked TeacherSlotPriorityValue = "blocked"
)

func FormTeacherSlotPriorityValues() []TeacherSlotPriorityValue {
	return []TeacherSlotPriorityValue{
		TeacherSlotBlocked,
	}
}

type TeacherSlotPriority struct {
	ID         uuid.UUID                `json:"id" db:"id"`
	TeacherID  uuid.UUID                `json:"teacher_id" db:"teacher_id"`
	TimeSlotID uuid.UUID                `json:"time_slot_id" db:"time_slot_id"`
	Priority   TeacherSlotPriorityValue `json:"priority" db:"priority"`
	CreatedAt  time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time                `json:"updated_at" db:"updated_at"`

	TimeSlot *LessonSlot `json:"time_slot,omitempty" db:"-"`
}

func (t TeacherSlotPriority) String() string {
	parts := make([]string, 0, 8)

	if t.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", t.ID))
	}

	parts = append(parts, fmt.Sprintf("teacher_id: %s", t.TeacherID))
	parts = append(parts, fmt.Sprintf("time_slot_id: %s", t.TimeSlotID))
	parts = append(parts, fmt.Sprintf("priority: %s", t.Priority))

	if !t.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", t.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", t.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("TeacherSlotPriority{%s}", strings.Join(parts, ", "))
}

func (t *TeacherSlotPriority) ValidatePriority() error {
	if slices.Index(FormTeacherSlotPriorityValues(), t.Priority) == -1 {
		return fmt.Errorf("invalid teacher slot priority value")
	}

	return nil
}

func (TeacherSlotPriority) TableName() string {
	return "teacher_slot_priorities"
}

func (TeacherSlotPriority) EntityName() string {
	return "teacher slot priority"
}
