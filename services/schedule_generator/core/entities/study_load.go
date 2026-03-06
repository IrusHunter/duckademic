package entities

import (
	"fmt"
	"slices"
)

// StudyLoad represents a university study load in the scheduling context.
//
// The model enforces workload constraints.
type StudyLoad struct {
	UnassignedLesson           // Base lesson definition.
	CurrentHours     int       // Currently scheduled hours.
	Lessons          []*Lesson // Scheduled lessons contributing to the load.
}

// NewStudyLoad creates a new StudyLoad instance.
//
// It requires an unassigned lesson entity (ul).
func NewStudyLoad(ul UnassignedLesson) *StudyLoad {
	return &StudyLoad{UnassignedLesson: ul}
}

// CheckLesson checks if the lesson can be added.
// It checks whether the lesson would exceed the remaining required hours.
//
// Return an error if validation fails.
func (sl *StudyLoad) CheckLesson(lesson *Lesson) error {
	if sl.IsEnoughHours() {
		return fmt.Errorf("enough hours for load")
	}

	return nil
}

// AddLesson registers the lesson.
//
// Uses CheckLesson for check.
func (sl *StudyLoad) AddLesson(lesson *Lesson) error {
	if err := sl.CheckLesson(lesson); err != nil {
		return fmt.Errorf("lesson check fails: %s", err.Error())
	}

	sl.CurrentHours += lesson.Value
	if sl.IsEnoughHours() {
		sl.CurrentHours += sl.CountHoursDeficit()
	}
	sl.Lessons = append(sl.Lessons, lesson)

	return nil
}

// CountHoursDeficit returns the number of missing study hours.
func (sl *StudyLoad) CountHoursDeficit() int {
	return sl.Discipline.GetRequiredHours(sl.Type) - sl.CurrentHours
}

// IsEnoughHours returns true if the study load has no remaining hours to schedule.
func (sl *StudyLoad) IsEnoughHours() bool {
	return sl.CountHoursDeficit() <= 0
}

// ==========================================================================================================
// =============================================== LoadService ==============================================
// ==========================================================================================================

// LoadService tracks and evaluates the study workload.
type LoadService interface {
	GetAssignedLessons() []*Lesson             // Returns registered lessons as an array.
	GetLessonTypes() []*LessonType             // Returns all lesson types from registered loads.
	AddLoad(*StudyLoad)                        // Registers a new study load.
	GetPreviousLessonOnDay(LessonSlot) *Lesson // Returns the nearest previous lesson on the slot's day.
}

// NewLoadService creates a new LoadService basic instance.
func NewLoadService() LoadService {
	return &loadService{}
}

// loadService is the basic implementation of the LoadService interface.
type loadService struct {
	loads []*StudyLoad
}

func (lc *loadService) GetAssignedLessons() (result []*Lesson) {
	for _, load := range lc.loads {
		result = append(result, load.Lessons...)
	}

	return result
}
func (lc *loadService) AddLoad(sl *StudyLoad) {
	lc.loads = append(lc.loads, sl)
}
func (lc *loadService) GetLessonTypes() (result []*LessonType) {
	for _, load := range lc.loads {
		ind := slices.IndexFunc(result, func(other *LessonType) bool {
			return load.Type == other
		})

		if ind == -1 {
			result = append(result, load.Type)
		}
	}

	return
}
func (lc *loadService) GetPreviousLessonOnDay(slot LessonSlot) *Lesson {
	lessons := lc.GetAssignedLessons()
	if len(lessons) == 0 {
		return nil
	}

	previous := lessons[0]
	for _, lesson := range lessons {
		if slot.After(lesson.LessonSlot) && lesson.After(previous) {
			previous = lesson
		}
	}

	if slot.After(previous.LessonSlot) && previous.Day == slot.Day {
		return previous
	}
	return nil
}
