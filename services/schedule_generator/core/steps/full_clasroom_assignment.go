package steps

import (
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/components"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/services"
)

type fullClassroomAssignmentStep struct {
	methods       map[components.ComponentIdentifier]components.ClassroomAssigner
	lessonService services.LessonService
	classrooms    []*entities.Classroom
}

func NewFullClassroomAssignmentStep(c *GeneratorContext) PipelineStep {
	s := fullClassroomAssignmentStep{methods: map[components.ComponentIdentifier]components.ClassroomAssigner{}}

	s.lessonService = c.fullData.lessonService
	s.classrooms = c.weekData.classroomService.GetAll()

	munkresClassroomAssigner := components.NewClassroomAssigner()
	s.methods[munkresClassroomAssigner.GetComponentIdentifier()] = munkresClassroomAssigner

	return &s
}

func (s *fullClassroomAssignmentStep) GetNextStep(c *GeneratorContext) PipelineStep {
	return nil
}
func (s *fullClassroomAssignmentStep) CanGoToTheNextStep() error {
	if lessonsWithoutC := s.lessonService.CountLessonsWithoutClassroom(); lessonsWithoutC != 0 {
		return fmt.Errorf("%d lessons haven't classrooms", lessonsWithoutC)
	}

	return nil
}
func (s *fullClassroomAssignmentStep) InsertData(data any) error {
	return nil
}
func (s *fullClassroomAssignmentStep) Process(cID components.ComponentIdentifier) (any, error) {
	comp, ok := s.methods[cID]
	if !ok {
		return nil, fmt.Errorf("component %q not allowed for weekly classroom assignment step", cID)
	}

	errs := comp.Run(components.ClassroomAssignerInput{
		Lessons:    s.lessonService.GetAll(),
		Classrooms: s.classrooms,
	})
	return responses.GeneratedLessonsWithC{
		LessonsWithClassroom:    responses.FormGeneratedLessons(s.lessonService.GetAll()),
		LessonsWithoutClassroom: errs,
	}, nil
}
