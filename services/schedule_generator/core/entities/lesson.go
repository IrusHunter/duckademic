package entities

import (
	"fmt"

	"github.com/google/uuid"
)

// Lesson represents an assigned lesson based on an UnsignedLesson.
type Lesson struct {
	ID         uuid.UUID
	*StudyLoad            // Base lesson data without time assignment.
	LessonSlot            // Assigned time slot
	Value      int        // Number of academic hours
	Classroom  *Classroom // Classroom where this lesson is scheduled
}

// NewLesson creates a new Lesson instance.
//
// It requires an unassigned lesson definition (ul),
// an assigned lesson slot (ls), and lesson value in academic hours (v).
func NewLesson(sl *StudyLoad, ls LessonSlot, v int) *Lesson {
	return &Lesson{
		StudyLoad:  sl,
		LessonSlot: ls,
		Value:      v,
		ID:         uuid.New(),
	}
}

// After returns true when other lesson (other) is positioned after the receiver.
// Otherwise returns false (even if slots are equal).
func (l *Lesson) After(other *Lesson) bool {
	return l.LessonSlot.After(other.LessonSlot)
}

// SetClassroom assigns the lesson to a selected classroom.
func (l *Lesson) SetClassroom(c *Classroom) error {
	if err := c.CheckLesson(l); err != nil {
		return fmt.Errorf("classroom %s unavailable for lesson: %w", c.RoomNumber, err)
	}

	l.Classroom = c

	if err := c.AddLesson(l); err != nil {
		panic("pass the check before, but error accurse")
	}

	return nil
}

// String returns a human-readable representation of the lesson.
//
// The output format is: "teacher: %%, student group: %%, discipline: %%, lesson type: %%, lesson slot: %%, classroom: %%".
func (l *Lesson) String() string {
	classroomStr := ""
	if l.Classroom != nil {
		classroomStr = fmt.Sprintf(", classroom: %s", l.Classroom.RoomNumber)
	}
	return fmt.Sprintf("teacher: %s, student group: %s, discipline: %s, lesson type: %s, lesson slot: %s%s",
		l.Teacher.UserName, l.StudentGroup.Name, l.Discipline.Name, l.Type.Name, l.LessonSlot.String(), classroomStr,
	)
}

// MoveLessonTo moves lesson to another slot (to).
func (l *Lesson) MoveLessonTo(to LessonSlot) {
	l.LessonSlot = to
}

// ==========================================================================================================
// =============================================== LessonSlot ===============================================
// ==========================================================================================================

// LessonSlot represents the position (coordinate) of a lesson within the schedule grid: Day is the day index,
// and Slot is the time-slot index within that day.
type LessonSlot struct {
	Day  int // Day position in the schedule grid.
	Slot int // Time slot position within the day.
}

// NewLessonSlot creates a new LessonSlot instance.
//
// It requires day (day) and slot (slot) for the new instance.
func NewLessonSlot(day, slot int) LessonSlot {
	return LessonSlot{
		Day:  day,
		Slot: slot,
	}
}

// String returns a human-readable representation of LessonSlot.
//
// The output format is: "day: %day, slot: %slot".
func (ls *LessonSlot) String() string {
	return fmt.Sprintf("day: %d, slot: %d", ls.Day, ls.Slot)
}

// After returns true when the receiver is positioned after other lesson slot (other) .
// Otherwise returns false (even if slots are equal).
func (ls *LessonSlot) After(other LessonSlot) bool {
	if ls.Day == other.Day {
		return ls.Slot > other.Slot
	}
	return ls.Day > other.Day
}

// ==========================================================================================================
// ============================================ UnassignedLesson ============================================
// ==========================================================================================================

// UnassignedLesson represents a lesson draft without confirmed assignments. Used for different errors.
type UnassignedLesson struct {
	Type         *LessonType   // Type of the lesson.
	Teacher      *Teacher      // Assigned teacher.
	StudentGroup *StudentGroup // Assigned group.
	Discipline   *Discipline   // Subject for the lesson.
}

// NewUnassignedLesson creates a new UnsignedLesson instance.
//
// It requires its own lesson type (lt), assigned teacher (t) and student group (sg) to it, and
// discipline (d) for the instance.
func NewUnassignedLesson(lt *LessonType, t *Teacher, sg *StudentGroup, d *Discipline) *UnassignedLesson {
	return &UnassignedLesson{
		Type:         lt,
		Teacher:      t,
		StudentGroup: sg,
		Discipline:   d,
	}
}

// Validate checks whether all required fields of UnassignedLesson are set.
// It returns an error describing the first missing field.
func (ul *UnassignedLesson) Validate() error {
	if ul.Type == nil {
		return fmt.Errorf("lesson type is not assigned")
	}
	if ul.Teacher == nil {
		return fmt.Errorf("teacher is not assigned")
	}
	if ul.StudentGroup == nil {
		return fmt.Errorf("student group is not assigned")
	}
	if ul.Discipline == nil {
		return fmt.Errorf("discipline is not assigned")
	}

	return nil
}
