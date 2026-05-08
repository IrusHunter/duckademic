package steps

import (
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/components"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/services"
)

type weekdayAllocationStep struct {
	methods             map[components.ComponentIdentifier]components.WeekdayAllocator
	studentGroupService services.StudentGroupService
}

func NewWeekdayAllocationStep(c *GeneratorContext) PipelineStep {
	s := weekdayAllocationStep{methods: map[components.ComponentIdentifier]components.WeekdayAllocator{}}

	evenWA := components.NewDayBlocker(c.fullData.studentGroupService.GetAll()[0].GetFullWeekCount(), c.config.LessonFillRate)
	s.methods[evenWA.GetComponentIdentifier()] = evenWA

	s.studentGroupService = c.weekData.studentGroupService

	return &s
}

func (s *weekdayAllocationStep) GetNextStep(c *GeneratorContext) PipelineStep {
	c.weekData.disciplineService.CutLoadTo(2)
	return NewWeeklyTimeSlotAssignmentStep(c)
}
func (s *weekdayAllocationStep) CanGoToTheNextStep() error {
	if count := s.studentGroupService.GetSlotDeficitOfReservedSlotsForLT(); count != 0 {
		return fmt.Errorf("student groups required %d additional reserved slots", count)
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

	errs := comp.Run(s.studentGroupService.GetAll())
	return responses.DaysForLessonTypes{
		StudentGroups: responses.FormDaysForLessonTypes(s.studentGroupService.GetAll()),
		Errors:        errs,
	}, nil
}
