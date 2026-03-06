package services

import (
	"fmt"
	"slices"
	"sort"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
)

// LessonService aggregates and manages lessons that the generator works with.
type LessonService interface {
	GetAll() []*entities.Lesson // Returns a slice with all lessons as pointers.
	// Assigns a lesson to the selected slot.
	AssignLesson(*entities.StudyLoad, entities.LessonSlot) error
	MoveLessonTo(*entities.Lesson, entities.LessonSlot) error // MoveLessonTo moves lesson to another slot (to).
	GetWeekLessons(int) []*entities.Lesson                    // TODO: collect bone lessons in another structure.
	// Sorts a slice of lessons according to the provided comparison function.
	Sort([]*entities.Lesson, func(a, b *entities.Lesson) int) []*entities.Lesson
	// Returns a comparison function that orders lessons by their lesson slot.
	ByLessonSlot(order int, next func(a, b *entities.Lesson) int) func(a, b *entities.Lesson) int
	// Serves as the final comparison in a sort chain, always returning 0.
	Equal(a, b *entities.Lesson) int
	CountLessonsWithoutClassroom([]*entities.Lesson) int // Returns the number of lessons that do not have an assigned classroom.
}

// NewLessonService creates a new LessonService basic instance.
//
// It requires a number of academic hours for lessons (lesson value - lv).
//
// Returns an error if the lesson value is below or equal to zero.
func NewLessonService(lv int) (LessonService, error) {
	if lv <= 0 {
		return nil, fmt.Errorf("lessonValue below/equal to 0 (%d)", lv)
	}

	ls := lessonService{lessonValue: lv}

	return &ls, nil
}

// lessonService is the basic implementation of the LessonService interface.
type lessonService struct {
	lessons     []*entities.Lesson
	lessonValue int
}

func (ls *lessonService) GetAll() []*entities.Lesson {
	return ls.lessons
}
func (ls *lessonService) AssignLesson(sl *entities.StudyLoad, slot entities.LessonSlot) error {
	lesson := entities.NewLesson(sl.UnassignedLesson, slot, ls.lessonValue)

	if err := sl.Teacher.CheckLesson(lesson); err != nil {
		return fmt.Errorf("lesson unavailable for teacher: %w", err)
	}
	if err := sl.StudentGroup.CheckLesson(lesson); err != nil {
		return fmt.Errorf("lesson unavailable for student group: %w", err)
	}
	if err := sl.CheckLesson(lesson); err != nil {
		return fmt.Errorf("lesson unavailable for study load: %w", err)
	}

	ls.lessons = append(ls.lessons, lesson)

	if err := lesson.StudentGroup.AddLesson(lesson); err != nil {
		panic("pass the check before, but error accurse")
	}
	if err := lesson.Teacher.AddLesson(lesson); err != nil {
		panic("pass the check before, but error accurse")
	}
	if err := sl.AddLesson(lesson); err != nil {
		panic("pass the check before, but error accurse")
	}

	return nil
}
func (ls *lessonService) GetWeekLessons(week int) (res []*entities.Lesson) {
	for _, l := range ls.lessons {
		if l.Day/7 == week {
			res = append(res, l)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i].Day != res[j].Day {
			return res[i].Day < res[j].Day
		}
		return res[i].Slot < res[j].Slot
	})
	return
}
func (ls *lessonService) MoveLessonTo(lesson *entities.Lesson, to entities.LessonSlot) error {
	if err := lesson.Teacher.LessonCanBeMoved(lesson.LessonSlot, to); err != nil {
		return err
	}
	if err := lesson.StudentGroup.LessonCanBeMoved(lesson, to); err != nil {
		return err
	}

	if err := lesson.Teacher.MoveLessonTo(lesson.LessonSlot, to); err != nil {
		panic("pass the check before, but error accurse")
	}
	if err := lesson.StudentGroup.MoveLessonTo(lesson, to); err != nil {
		panic("pass the check before, but error accurse")
	}
	lesson.MoveLessonTo(to)
	return nil
}
func (ls *lessonService) Sort(lessons []*entities.Lesson, sortFunc func(a, b *entities.Lesson) int) []*entities.Lesson {
	result := make([]*entities.Lesson, len(lessons))
	copy(result, lessons)

	slices.SortFunc(result, sortFunc)

	return result
}
func (ls *lessonService) ByLessonSlot(order int, next func(a, b *entities.Lesson) int) func(a, b *entities.Lesson) int {
	return func(a, b *entities.Lesson) int {
		if a.After(b) {
			return 1 * order
		} else if b.After(a) {
			return -1 * order
		}

		return next(a, b)
	}
}
func (ls *lessonService) Equal(a, b *entities.Lesson) int {
	return 0
}
func (ls *lessonService) CountLessonsWithoutClassroom(l []*entities.Lesson) (result int) {
	for _, lesson := range l {
		if lesson.Classroom == nil {
			result++
		}
	}
	return
}

// func (ls *lessonService) Select(filterFunc func(a *entities.Lesson) bool) []*entities.Lesson {
// 	result := make([]*entities.Lesson, 0)
//
// 	for _, lesson := range ls.lessons {
// 		if filterFunc(lesson) {
// 			result = append(result, lesson)
// 		}
// 	}
//
// 	return result
// }
