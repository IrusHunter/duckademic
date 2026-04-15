package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type LessonOccurrenceStatus string

const (
	LessonOccurrenceScheduled LessonOccurrenceStatus = "scheduled"
	LessonOccurrenceCanceled  LessonOccurrenceStatus = "canceled"
	LessonOccurrenceCompleted LessonOccurrenceStatus = "completed"
)

type LessonOccurrence struct {
	ID             uuid.UUID              `json:"id" db:"id"`
	StudyLoadID    uuid.UUID              `json:"study_load_id" db:"study_load_id"`
	TeacherID      uuid.UUID              `json:"teacher_id" db:"teacher_id"`
	StudentGroupID uuid.UUID              `json:"student_group_id" db:"student_group_id"`
	LessonSlotID   uuid.UUID              `json:"lesson_slot_id" db:"lesson_slot_id"`
	Date           time.Time              `json:"date" db:"date"`
	ClassroomID    *uuid.UUID             `json:"classroom_id,omitempty" db:"classroom_id"`
	Status         LessonOccurrenceStatus `json:"status" db:"status"`
	MovedToID      *uuid.UUID             `json:"moved_to_id,omitempty" db:"moved_to_id"`
	MovedFromID    *uuid.UUID             `json:"moved_from_id,omitempty" db:"moved_from_id"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" db:"updated_at"`
}

func (l LessonOccurrence) String() string {
	parts := make([]string, 0, 12)

	if l.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", l.ID))
	}

	parts = append(parts, fmt.Sprintf("study_load_id: %s", l.StudyLoadID))
	parts = append(parts, fmt.Sprintf("teacher_id: %s", l.TeacherID))
	parts = append(parts, fmt.Sprintf("student_group_id: %s", l.StudentGroupID))
	parts = append(parts, fmt.Sprintf("lesson_slot_id: %s", l.LessonSlotID))
	parts = append(parts, fmt.Sprintf("date: %s", l.Date.Format(db.TimeFormat)))
	parts = append(parts, fmt.Sprintf("status: %s", l.Status))

	if l.ClassroomID != nil {
		parts = append(parts, fmt.Sprintf("classroom_id: %s", *l.ClassroomID))
	}
	if l.MovedToID != nil {
		parts = append(parts, fmt.Sprintf("moved_to_id: %s", *l.MovedToID))
	}
	if l.MovedFromID != nil {
		parts = append(parts, fmt.Sprintf("moved_from_id: %s", *l.MovedFromID))
	}

	if !l.CreatedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("created_at: %s", l.CreatedAt.Format(db.TimeFormat)))
		parts = append(parts, fmt.Sprintf("updated_at: %s", l.UpdatedAt.Format(db.TimeFormat)))
	}

	return fmt.Sprintf("LessonOccurrence{%s}", strings.Join(parts, ", "))
}

func (LessonOccurrence) TableName() string {
	return "lesson_occurrences"
}

func (LessonOccurrence) EntityName() string {
	return "lesson occurrence"
}
