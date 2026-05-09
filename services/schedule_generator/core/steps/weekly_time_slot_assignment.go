package steps

import (
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/components"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/services"
)

type weeklyTimeSlotAssignmentStep struct {
	methods       map[components.ComponentIdentifier]components.TimeSlotAssigner
	studyLoads    services.StudyLoadService
	lessonService services.LessonService
}

func NewWeeklyTimeSlotAssignmentStep(c *GeneratorContext) PipelineStep {
	s := weeklyTimeSlotAssignmentStep{methods: map[components.ComponentIdentifier]components.TimeSlotAssigner{}}

	s.studyLoads, _ = services.NewStudyLoadService(c.weekData.studyLoadService.GetAll())
	s.lessonService = c.weekData.lessonService

	onePerWeekTimeSlotAssigner := components.NewBoneGenerator()
	s.methods[onePerWeekTimeSlotAssigner.GetComponentIdentifier()] = onePerWeekTimeSlotAssigner

	bruteTimeSlotAssigner := components.NewMissingLessonAdder()
	s.methods[bruteTimeSlotAssigner.GetComponentIdentifier()] = bruteTimeSlotAssigner

	return &s
}

func (s *weeklyTimeSlotAssignmentStep) GetNextStep(c *GeneratorContext) PipelineStep {
	return NewWeeklyClassroomAssignmentStep(c)
}
func (s *weeklyTimeSlotAssignmentStep) CanGoToTheNextStep() error {
	if _, hoursDept := s.studyLoads.GetHoursDeficit(); hoursDept != 0 {
		return fmt.Errorf("%d study loads haven't assigned week lesson", hoursDept)
	}

	return nil
}
func (s *weeklyTimeSlotAssignmentStep) InsertData(data any) error {
	return nil
}
func (s *weeklyTimeSlotAssignmentStep) Process(cID components.ComponentIdentifier) (any, error) {
	comp, ok := s.methods[cID]
	if !ok {
		return nil, fmt.Errorf("component %q not allowed for weekly time slot assignment step", cID)
	}

	errs := comp.Run(components.TimeSlotAssignerInput{
		StudyLoads:    s.studyLoads.GetAll(),
		LessonService: s.lessonService,
	})
	return responses.GeneratedLessons{
		Lessons: responses.FormGeneratedLessons(s.lessonService.GetAll()),
		Errors:  errs,
	}, nil
}
