package steps

import (
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/components"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
)

type weekdayAllocationStep struct {
	methods            map[components.ComponentIdentifier]components.WeekdayAllocator
	studentGroups      []*entities.StudentGroup
	canGoToTheNextStep bool
}

func NewWeekdayAllocationStep(c *GeneratorContext) PipelineStep {
	s := weekdayAllocationStep{methods: map[components.ComponentIdentifier]components.WeekdayAllocator{}}

	evenWA := components.NewDayBlocker(c.fullData.studentGroupService.GetAll()[0].GetFullWeekCount(), c.config.LessonFillRate)
	s.methods[evenWA.GetComponentIdentifier()] = evenWA

	s.studentGroups = c.weekData.studentGroupService.GetAll()

	return &s
}

func (s *weekdayAllocationStep) GetNextStep(c *GeneratorContext) PipelineStep {
	c.weekData.disciplineService.CutLoadTo(2)
	return NewWeeklyTimeSlotAssignmentStep(c)
}
func (s *weekdayAllocationStep) CanGoToTheNextStep() error {
	if !s.canGoToTheNextStep {
		return fmt.Errorf("required action was not applied")
	}
	return nil
}
func (s *weekdayAllocationStep) InsertData(data any) error {
	return nil
}
func (s *weekdayAllocationStep) Process(cID components.ComponentIdentifier) (any, error) {
	comp, ok := s.methods[cID]
	if !ok {
		return responses.DaysForLessonTypes{},
			fmt.Errorf("component %q not allowed for weekday allocation step", cID)
	}

	s.canGoToTheNextStep = true
	errs := comp.Run(s.studentGroups)
	return responses.DaysForLessonTypes{
		StudentGroups: responses.FormDaysForLessonTypes(s.studentGroups),
		Errors:        errs,
	}, nil
}
