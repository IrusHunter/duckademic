package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type LessonSlot struct {
	ID        uuid.UUID     `json:"id" db:"id"`
	Slot      int           `json:"slot" db:"slot"`
	Weekday   time.Weekday  `json:"weekday" db:"weekday"`
	StartTime time.Duration `json:"start_time" db:"start_time"`
	Duration  time.Duration `json:"duration" db:"duration"`
	CreatedAt time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" json:"updated_at"`
}

func (l LessonSlot) String() string {
	parts := make([]string, 0, 6)

	if l.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", l.ID))
	}

	parts = append(parts, fmt.Sprintf("slot: %d", l.Slot))
	parts = append(parts, fmt.Sprintf("weekday: %s", l.Weekday.String()))
	parts = append(parts, fmt.Sprintf("start_time: %s", l.StartTime.String()))
	parts = append(parts, fmt.Sprintf("duration: %s", l.Duration.String()))

	if !l.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", l.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", l.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("LessonSlot{%s}", strings.Join(parts, ", "))
}

func (LessonSlot) TableName() string {
	return "lesson_slots"
}
