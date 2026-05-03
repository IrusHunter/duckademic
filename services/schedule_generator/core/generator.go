package core

import (
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/components"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/steps"
	externalEntities "github.com/IrusHunter/duckademic/services/schedule_generator/entities"
)

type GeneratorStep string

const (
	Setup                               GeneratorStep = "SETUP"
	DayBlocking                         GeneratorStep = "DAY_BLOCKING"
	BoneLessonBuilding                  GeneratorStep = "BONE_LESSON_BUILDING"
	ToBoneLessonsClassroomAssigning     GeneratorStep = "TO_BONE_LESSONS_CLASSROOM_ASSIGNING"
	LessonSkeletonBuilding              GeneratorStep = "LESSON_SKELETON_BUILDING"
	FloatingLessonAdding                GeneratorStep = "FLOATING_LESSON_ADDING"
	ToFloatingLessonsClassroomAssigning GeneratorStep = "TO_FLOATING_LESSONS_CLASSROOM_ASSIGNING"
	Extraction                          GeneratorStep = "EXTRACTION"
)

type ScheduleGenerator struct {
	steps.GeneratorContext
	currentStep steps.PipelineStep
}

func NewScheduleGenerator(cfg externalEntities.ScheduleGeneratorConfig) (*ScheduleGenerator, error) {
	c, err := steps.NewGeneratorContext(cfg)
	if err != nil {
		return nil, err
	}

	return &ScheduleGenerator{GeneratorContext: *c}, nil
}

func (g *ScheduleGenerator) SubmitAndGoToTheNextStep(ignoreWarnings bool) error {
	if g.currentStep == nil {
		if !g.AllDataReceived() {
			return fmt.Errorf("not all data received")
		}

		g.currentStep = steps.NewWeekdayAllocationStep(&g.GeneratorContext)
		return nil
	}

	err := g.currentStep.CanGoToTheNextStep()
	if err != nil && !ignoreWarnings {
		return err
	}

	nextStep := g.currentStep.GetNextStep(&g.GeneratorContext)
	if nextStep == nil {
		return nil
		// if err == nil {
		// 	return components.NewUnexpectedError("failed to go to the next step with: no error from current step",
		// 		"Generator", "SubmitAndGoToTheNextStep", fmt.Errorf("no error from current step but should be"))
		// }
		// return err
	}

	g.currentStep = nextStep
	return nil
}

func (g *ScheduleGenerator) ProcessStep(method components.ComponentIdentifier, data any) (any, error) {
	if err := g.currentStep.InsertData(data); err != nil {
		return nil, err
	}

	return g.currentStep.Process(method)
}
