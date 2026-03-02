package entities

import (
	"fmt"

	"github.com/google/uuid"
)

// Classroom represents a university classroom in a scheduling context.
//
// The model disallows simultaneous classes and prevents capacity overflow.
type Classroom struct {
	BusyGrid                 // Availability grid.
	ID             uuid.UUID // Unique identifier of the Classroom.
	RoomNumber     string    // Human-readable identifier of the Classroom.
	Capacity       int       // Maximum number of students allowed in the classroom.
	FillPercentage float32   // Capacity fraction above which the classroom is full.
	lessons        []*Lesson // Lessons scheduled in this classroom.
}

// NewClassroom create a new Classroom instance.
//
// It requires the classroom id, room number (rn), capacity of the classroom,
// busy grid (bg), and percentage of used capacity (fp).
func NewClassroom(id uuid.UUID, rn string, c int, bg *BusyGrid, fp float32) *Classroom {
	return &Classroom{
		ID:             id,
		RoomNumber:     rn,
		Capacity:       c,
		BusyGrid:       *bg,
		FillPercentage: fp,
	}
}

// CanAccommodate returns true if the classroom has enough capacity. Otherwise returns false.
func (c *Classroom) CanAccommodate(number int) bool {
	return float32(c.Capacity)*c.FillPercentage >= float32(number)
}

// CheckLesson checks if a lesson can be assigned to the classroom. It checks that the classroom has
// enough capacity and availability.
//
// Return an error if validation fails.
func (c *Classroom) CheckLesson(lesson *Lesson) error {
	if !c.IsFree(lesson.LessonSlot) {
		return fmt.Errorf("slot %s is busy", lesson.LessonSlot.String())
	}

	if !c.CanAccommodate(lesson.StudentGroup.StudentNumber) {
		return fmt.Errorf("can't accommodate %d people", lesson.StudentGroup.StudentNumber)
	}

	return nil
}

// AddLesson registers the lesson at all dependent services.
//
// Returns an error if CheckLesson fails.
func (c *Classroom) AddLesson(lesson *Lesson) error {
	if err := c.CheckLesson(lesson); err != nil {
		return fmt.Errorf("lesson check fails: %w", err)
	}

	c.lessons = append(c.lessons, lesson)
	c.BusyGrid.SetSlotBusyState(lesson.LessonSlot, true)

	return nil
}

// CountOverflowLessons returns the number of lessons that exceed the classroom capacity.
func (c *Classroom) CountOverflowLessons() (result int) {
	for _, lesson := range c.lessons {
		if !c.CanAccommodate(lesson.StudentGroup.StudentNumber) {
			result++
		}
	}
	return
}

// CountLessonOverlapping returns the number of overlapping lessons. Counts only lessons that overlap.
func (c *Classroom) CountLessonOverlapping() int {
	return c.BusyGrid.CountLessonOverlapping(c.lessons)
}
