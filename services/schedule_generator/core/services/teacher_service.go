package services

import (
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
	"github.com/google/uuid"
)

// TeacherService aggregates and manages teachers that the generator works with.
type TeacherService interface {
	// Returns a pointer to the teacher with the given ID.
	Find(uuid.UUID) *entities.Teacher
	// Returns a slice with all teachers as pointers.
	GetAll() []*entities.Teacher
	// Returns windows (gaps between busy slots) and the sum of it.
	GetWindows() ([]responses.TeacherWindow, int)
	// Returns the overlapping lesson slots and the sum of it.
	GetLessonOverlapping() ([]responses.TeacherLessonOverlap, int)
}

// NewTeacherService creates a new TeacherService instance.
//
// It requires an array of teachers (t).
func NewTeacherService(t []*entities.Teacher) TeacherService {
	return &teacherService{teachers: t}
}

// teacherService is the basic implementation of the TeacherService interface.
type teacherService struct {
	teachers []*entities.Teacher
}

func (ts *teacherService) GetAll() []*entities.Teacher {
	return ts.teachers
}
func (ts *teacherService) Find(id uuid.UUID) *entities.Teacher {
	for i := range ts.teachers {
		if ts.teachers[i].ID == id {
			return ts.teachers[i]
		}
	}

	return nil
}
func (ts *teacherService) GetWindows() (res []responses.TeacherWindow, count int) {
	res = []responses.TeacherWindow{}
	for _, t := range ts.teachers {
		windows := t.GetWindows()
		if len(windows) != 0 {
			count += len(windows)
			res = append(res, responses.TeacherWindow{
				CommonEntity: responses.CommonEntity{
					ID:   t.ID,
					Name: t.UserName,
				},
				Windows: responses.FormTimeSlots(windows),
			})
		}
	}
	return
}
func (ts *teacherService) GetLessonOverlapping() (res []responses.TeacherLessonOverlap, count int) {
	res = []responses.TeacherLessonOverlap{}

	for _, teacher := range ts.teachers {
		lo := teacher.GetLessonOverlapping()
		if len(lo) != 0 {
			count += len(lo)
			res = append(res, responses.TeacherLessonOverlap{
				CommonEntity: responses.CommonEntity{
					ID:   teacher.ID,
					Name: teacher.UserName,
				},
				OverlapLessons: responses.FormTimeSlots(lo),
			})
		}
	}

	return
}
