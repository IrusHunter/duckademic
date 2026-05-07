package services

import (
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
	externalEntities "github.com/IrusHunter/duckademic/services/schedule_generator/entities"
	"github.com/google/uuid"
)

// ClassroomService aggregates and manages classrooms that the generator works with.
type ClassroomService interface {
	// Returns a slice with all classrooms as pointers.
	GetAll() []*entities.Classroom
	// Returns a pointer to the classroom with the given ID.
	Find(uuid.UUID) *entities.Classroom
	// Returns the lessons that exceed the classrooms capacity and the sum of them.
	GetOverflowLessons() ([]responses.ClassroomOverflow, int)
	// Returns the overlapping lessons and the sum of it.
	GetLessonOverlapping() ([]responses.ClassroomLessonOverlap, int)
}

// NewClassroomService creates a new ClassroomService basic instance.
//
// It requires the array of database classrooms (c), the busy grid for them (bg),
// and percentage of used capacity (fp).
//
// Returns an error if any classroom is an invalid model.
func NewClassroomService(c []externalEntities.Classroom, bg [][]float32, fp float32) (ClassroomService, error) {
	classrooms := make([]*entities.Classroom, len(c))

	for i := range c {
		classrooms[i] = entities.NewClassroom(
			c[i].ID, c[i].Number, c[i].Capacity, entities.NewBusyGrid(bg), fp,
		)
	}

	return &classroomService{classrooms: classrooms}, nil
}

// classroomService is the basic implementation of the ClassroomService interface.
type classroomService struct {
	classrooms []*entities.Classroom
}

func (s *classroomService) GetAll() []*entities.Classroom {
	return s.classrooms
}
func (s *classroomService) Find(classroomID uuid.UUID) *entities.Classroom {
	for i := range s.classrooms {
		if s.classrooms[i].ID == classroomID {
			return s.classrooms[i]
		}
	}

	return nil
}
func (s *classroomService) GetOverflowLessons() (res []responses.ClassroomOverflow, count int) {
	res = []responses.ClassroomOverflow{}

	for _, classroom := range s.classrooms {
		l := classroom.GetOverflowLessons()
		if len(l) != 0 {
			count += len(l)
			res = append(res, responses.ClassroomOverflow{
				CommonEntity: responses.FormCommonEntity(classroom.ID, classroom.RoomNumber),
				Overflows:    responses.FormClassroomOverflowItems(l),
				Capacity:     classroom.GetMaximumCapacity(),
			})
		}
	}

	return
}
func (s *classroomService) GetLessonOverlapping() (res []responses.ClassroomLessonOverlap, count int) {
	res = []responses.ClassroomLessonOverlap{}

	for _, classroom := range s.classrooms {
		lo := classroom.GetLessonOverlapping()
		if len(lo) != 0 {
			count += len(lo)
			res = append(res, responses.ClassroomLessonOverlap{
				CommonEntity: responses.CommonEntity{
					ID:   classroom.ID,
					Name: classroom.RoomNumber,
				},
				OverlapLessons: responses.FormTimeSlots(lo),
			})
		}
	}

	return
}
