package services

import (
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	externalEntities "github.com/IrusHunter/duckademic/services/schedule_generator/entities"
)

// ClassroomService aggregates and manages classrooms that the generator works with.
type ClassroomService interface {
	GetAll() []*entities.Classroom // Returns a slice with all classrooms as pointers.
	CountOverflowLessons() int     // Returns the number of lessons that exceed the classrooms capacity.
	CountLessonOverlapping() int   // Returns the count of overlapping lessons.
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
func (s *classroomService) CountOverflowLessons() (result int) {
	for _, classroom := range s.classrooms {
		result += classroom.CountOverflowLessons()
	}
	return
}
func (s *classroomService) CountLessonOverlapping() (result int) {
	for _, classroom := range s.classrooms {
		result += classroom.CountLessonOverlapping()
	}
	return
}
