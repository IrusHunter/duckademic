package entities

import (
	"fmt"
	"slices"

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
	Lessons        []*Lesson // Lessons scheduled in this classroom.
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

func (c *Classroom) GetMaximumCapacity() int {
	return int(float32(c.Capacity) * c.FillPercentage)
}

// CanAccommodate returns true if the classroom has enough capacity. Otherwise returns false.
func (c *Classroom) CanAccommodate(number int) bool {
	return c.GetMaximumCapacity() >= number
}

// CheckLesson checks if a lesson can be assigned to the classroom. It checks that the classroom has
// enough capacity and availability.
//
// Return an error if validation fails.
func (c *Classroom) CheckLesson(lesson *Lesson) error {
	if !c.IsFree(lesson.LessonSlot) {
		return fmt.Errorf("%s is busy", lesson.LessonSlot.String())
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

	c.Lessons = append(c.Lessons, lesson)
	c.BusyGrid.SetSlotBusyState(lesson.LessonSlot, true)

	return nil
}

// COMMENT THIS
func (c *Classroom) RemoveLesson(lesson *Lesson) error {
	if lesson.Classroom != c {
		return fmt.Errorf("invalid lesson for classroom")
	}

	ind := slices.Index(c.Lessons, lesson)
	if ind == -1 {
		return fmt.Errorf("lesson not found for classroom")
	}
	c.Lessons = append(c.Lessons[:ind], c.Lessons[ind+1:]...)
	c.BusyGrid.SetSlotBusyState(lesson.LessonSlot, false)

	return nil
}

// COMMENT THIS
func (c *Classroom) GetOverflowLessons() []*Lesson {
	res := []*Lesson{}

	for _, lesson := range c.Lessons {
		if !c.CanAccommodate(lesson.StudentGroup.StudentNumber) {
			res = append(res, lesson)
		}
	}

	return res
}

// COMMENT THIS
func (c *Classroom) GetLessonOverlapping() []LessonSlot {
	return c.BusyGrid.GetLessonOverlapping(c.Lessons)
}
