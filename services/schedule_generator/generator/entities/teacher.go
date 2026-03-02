package entities

import (
	"fmt"

	"github.com/google/uuid"
)

// Teacher represents a university teacher in the scheduling context.
//
// The model enforces disallows simultaneous classes.
//
// TODO: add teacher availability constraints.
type Teacher struct {
	BusyGrid              // Availability grid.
	LoadService           // Manages and validates assigned study loads.
	ID          uuid.UUID // Unique identifier of the Teacher.
	UserName    string    // Human-readable identifier of the Teacher.
	Priority    int       // Higher value means higher priority (used for sorting).
}

// NewTeacher creates a new Teacher instance.
//
// It requires teacher's id, name (un), priority (p), busy grid for teacher (bg), and load service (ls).
func NewTeacher(id uuid.UUID, un string, p int, bg *BusyGrid, ls LoadService) *Teacher {
	return &Teacher{
		BusyGrid:    *bg,
		ID:          id,
		UserName:    un,
		Priority:    p,
		LoadService: ls,
	}
}

// NewTeacher creates a new Teacher instance with default configuration.
//
// It requires teacher's id, name (un), priority (p) and busy grid for teacher (bg).
func NewDefaultTeacher(id uuid.UUID, un string, p int, bg *BusyGrid) *Teacher {
	return NewTeacher(id, un, p, bg, NewLoadService())
}

// AddLesson register the lesson.
//
// Uses CheckLesson for check.
func (t *Teacher) AddLesson(lesson *Lesson) error {
	err := t.CheckLesson(lesson)
	if err != nil {
		return err
	}

	t.SetSlotBusyState(lesson.LessonSlot, true)

	return err
}

// CheckLesson checks if the lesson can be added. It checks slot validation and availability.
//
// Return an error if validation fails.
func (t *Teacher) CheckLesson(lesson *Lesson) error {
	if err := t.CheckSlot(lesson.LessonSlot); err != nil {
		return err
	}
	if !t.IsFree(lesson.LessonSlot) {
		return fmt.Errorf("teacher is busy")
	}

	return nil
}

// CountLessonOverlapping returns the count of overlapping lessons. Counts only lessons that overlap.
func (t *Teacher) CountLessonOverlapping() int {
	return t.BusyGrid.CountLessonOverlapping(t.GetAssignedLessons())
}

// CalculateClassValueFor returns own value for a classroom.
// It factors in the classroom of the previous lesson on that day.
func (t *Teacher) CalculateClassValueFor(l *Lesson, c *Classroom) float32 {
	previousLesson := t.GetPreviousLessonOnDay(l.LessonSlot)
	if previousLesson == nil {
		return 0.5
	}

	if previousLesson.Classroom == c {
		return 0.5
	}
	return 1
}
