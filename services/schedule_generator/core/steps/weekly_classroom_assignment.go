package steps

import (
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/components"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/services"
)

type weeklyClassroomAssignmentStep struct {
	methods       map[components.ComponentIdentifier]components.ClassroomAssigner
	lessonService services.LessonService
	classrooms    []*entities.Classroom
}

func NewWeeklyClassroomAssignmentStep(c *GeneratorContext) PipelineStep {
	s := weeklyClassroomAssignmentStep{methods: map[components.ComponentIdentifier]components.ClassroomAssigner{}}

	s.lessonService = c.weekData.lessonService
	s.classrooms = c.weekData.classroomService.GetAll()

	munkresClassroomAssigner := components.NewClassroomAssigner()
	s.methods[munkresClassroomAssigner.GetComponentIdentifier()] = munkresClassroomAssigner

	return &s
}

func (s *weeklyClassroomAssignmentStep) GetNextStep(c *GeneratorContext) PipelineStep {
	return NewWeeklyScheduleExpansionStep(c)
}
func (s *weeklyClassroomAssignmentStep) CanGoToTheNextStep() error {
	if lessonsWithoutC := len(s.lessonService.GetLessonsWithoutClassroom()); lessonsWithoutC != 0 {
		return fmt.Errorf("%d lessons haven't classrooms", lessonsWithoutC)
	}

	return nil
}
func (s *weeklyClassroomAssignmentStep) InsertData(data any) error {
	return nil
}
func (s *weeklyClassroomAssignmentStep) Process(cID components.ComponentIdentifier) (any, error) {
	comp, ok := s.methods[cID]
	if !ok {
		return nil, fmt.Errorf("component %q not allowed for weekly classroom assignment step", cID)
	}

	errs := comp.Run(components.ClassroomAssignerInput{
		Lessons:    s.lessonService.GetAll(),
		Classrooms: s.classrooms,
	})
	return responses.BoneLessonsWithC{
		LessonsWithClassroom:    responses.FormBoneLessons(s.lessonService.GetAll()),
		LessonsWithoutClassroom: errs,
	}, nil
}
