package services

import (
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
	externalEntities "github.com/IrusHunter/duckademic/services/schedule_generator/entities"
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
// It requires an array of database teachers (t) and a busy grid for them (bg).
//
// Returns an error if any teacher is an invalid model.
func NewTeacherService(t []externalEntities.Teacher, bg [][]float32) (TeacherService, error) {
	ts := teacherService{teachers: make([]*entities.Teacher, 0, len(t))}

	for i := range t {
		teacher := entities.NewDefaultTeacher(t[i].ID, t[i].Name, t[i].Priority, entities.NewBusyGrid(bg))
		// for _, day := range t[i].BusyDays {
		// 	err := teacher.BlockWeekDay(int(day))
		// 	if err != nil {
		// 		return nil, fmt.Errorf("teacher %s (%s) has invalid busy day %d (err: %s)",
		// 			teacher.UserName, teacher.ID, day, err.Error(),
		// 		)
		// 	}
		// }

		// sort in priority order (not necessary now)
		// success := false
		// for j, lowerTeacher := range ts.teachers {
		// 	if lowerTeacher.Priority <= teacher.Priority {
		// 		ts.teachers = append(ts.teachers[:j], append([]*entities.Teacher{teacher}, ts.teachers[j:]...)...)
		// 		success = true
		// 		break
		// 	}
		// }
		// if !success {
		ts.teachers = append(ts.teachers, teacher)
		// }
	}

	return &ts, nil
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
