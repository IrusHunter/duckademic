package services

import (
	"fmt"
	"slices"
	"sort"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
)

// LessonService aggregates and manages lessons that the generator works with.
type LessonService interface {
	// Returns a slice with all lessons as pointers.
	GetAll() []*entities.Lesson
	// Assigns a lesson to the selected slot.
	AssignLesson(*entities.StudyLoad, entities.LessonSlot) error
	// MoveLessonTo moves lesson to another slot (to).
	MoveLessonTo(*entities.Lesson, entities.LessonSlot) error
	// TODO: collect bone lessons in another structure.
	GetWeekLessons(int) []*entities.Lesson
	// Returns the number of lessons that do not have an assigned classroom.
	CountLessonsWithoutClassroom([]*entities.Lesson) int
	Select() *LessonSelector
}

// NewLessonService creates a new LessonService basic instance.
//
// Returns an error if the lesson value is below or equal to zero.
func NewLessonService() (LessonService, error) {
	return &lessonService{}, nil
}

// lessonService is the basic implementation of the LessonService interface.
type lessonService struct {
	lessons []*entities.Lesson
}

func (ls *lessonService) GetAll() []*entities.Lesson {
	return ls.lessons
}
func (ls *lessonService) AssignLesson(sl *entities.StudyLoad, slot entities.LessonSlot) error {
	lesson := entities.NewLesson(sl.UnassignedLesson, slot, sl.Type.Value)

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
func (ls *lessonService) CountLessonsWithoutClassroom(l []*entities.Lesson) (result int) {
	for _, lesson := range l {
		if lesson.Classroom == nil {
			result++
		}
	}
	return
}
func (ls *lessonService) Select() *LessonSelector {
	return &LessonSelector{
		lessons: ls.lessons,
	}
}

type LessonSelector struct {
	lessons []*entities.Lesson
}

func (ls *LessonSelector) Sort() *LessonSorter {
	return &LessonSorter{
		lessons:    ls.lessons,
		comparator: func(a, b *entities.Lesson) int { return 0 },
	}
}

type LessonSorter struct {
	lessons    []*entities.Lesson
	comparator func(a, b *entities.Lesson) int
}

func (s *LessonSorter) ByLessonSlot(order int) *LessonSorter {
	prev := s.comparator

	s.comparator = func(a, b *entities.Lesson) int {
		if a.After(b) {
			return 1 * order
		} else if b.After(a) {
			return -1 * order
		}
		return prev(a, b)
	}

	return s
}
func (s *LessonSorter) ToSlice() []*entities.Lesson {
	result := make([]*entities.Lesson, len(s.lessons))
	copy(result, s.lessons)

	slices.SortFunc(result, s.comparator)

	return result
}
func (s *LessonSorter) Last() *entities.Lesson {
	return s.lessons[len(s.lessons)-1]
}
