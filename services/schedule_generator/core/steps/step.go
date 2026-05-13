package steps

import (
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/components"
	"github.com/google/uuid"
)

type PipelineStep interface {
	GetNextStep(*GeneratorContext) PipelineStep
	CanGoToTheNextStep() error
	InsertData(data any) error
	Process(components.ComponentIdentifier) (any, error)
	ApplyManualChange(data map[string]string) error
}

type WeekdayBindingOverride struct {
	StudentGroupID uuid.UUID
	LessonTypeID   uuid.UUID
	OldWeekday     *int
	NewWeekday     *int
}

type TimeSlotForLessonOverride struct {
	StudyLoadID uuid.UUID
	OldDay      *int
	OldSlot     *int
	NewDay      *int
	NewSlot     *int
}

type ClassroomForLessonOverride struct {
	LessonID       uuid.UUID
	NewClassroomID *uuid.UUID
}
