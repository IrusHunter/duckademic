package steps

import (
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/components"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/services"
)

type fullTimeSlotAssignmentStep struct {
	methods          map[components.ComponentIdentifier]components.TimeSlotAssigner
	studyLoadService services.StudyLoadService
	lessonService    services.LessonService
}

func NewFullTimeSlotAssignmentStep(c *GeneratorContext) PipelineStep {
	s := fullTimeSlotAssignmentStep{methods: map[components.ComponentIdentifier]components.TimeSlotAssigner{}}

	s.studyLoadService, _ = services.NewStudyLoadService(c.fullData.studyLoadService.GetAll())
	s.lessonService = c.floatLessonService

	onePerWeekTimeSlotAssigner := components.NewBoneGenerator()
	s.methods[onePerWeekTimeSlotAssigner.GetComponentIdentifier()] = onePerWeekTimeSlotAssigner

	bruteTimeSlotAssigner := components.NewMissingLessonAdder()
	s.methods[bruteTimeSlotAssigner.GetComponentIdentifier()] = bruteTimeSlotAssigner

	return &s
}

func (s *fullTimeSlotAssignmentStep) GetNextStep(c *GeneratorContext) PipelineStep {
	return NewFullClassroomAssignmentStep(c)
}
func (s *fullTimeSlotAssignmentStep) CanGoToTheNextStep() error {
	if _, hoursDept := s.studyLoadService.GetHoursDeficit(); hoursDept != 0 {
		return fmt.Errorf("%d study loads haven't assigned week lesson", hoursDept)
	}

	return nil
}
func (s *fullTimeSlotAssignmentStep) InsertData(data any) error {
	return nil
}
func (s *fullTimeSlotAssignmentStep) Process(cID components.ComponentIdentifier) (any, error) {
	comp, ok := s.methods[cID]
	if !ok {
		return nil, fmt.Errorf("component %q not allowed for full time slot assignment step", cID)
	}

	errs := comp.Run(components.TimeSlotAssignerInput{
		StudyLoads:    s.studyLoadService.GetAll(),
		LessonService: s.lessonService,
	})
	return responses.GeneratedLessons{
		Lessons: responses.FormGeneratedLessons(s.lessonService.GetAll()),
		Errors:  errs,
	}, nil
}
