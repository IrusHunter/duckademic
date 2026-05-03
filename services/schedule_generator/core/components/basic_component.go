package components

import (
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/services"
)

// GeneratorComponent represents any component participating in the generation process.
// TI - input type of the component, TR - response type.
type GeneratorComponent[TI, TR any] interface {
	GetComponentIdentifier() ComponentIdentifier
	Run(TI) TR // The main improvement of schedule for generator
}

type WeekdayAllocator interface {
	GeneratorComponent[[]*entities.StudentGroup, []responses.LessonTypeDayDebt]
}

type TimeSlotAssignerInput struct {
	StudyLoads    []*entities.StudyLoad
	LessonService services.LessonService
}

type TimeSlotAssigner interface {
	GeneratorComponent[TimeSlotAssignerInput, []responses.UnassignedLesson]
}

type ClassroomAssignerInput struct {
	Lessons    []*entities.Lesson
	Classrooms []*entities.Classroom
}

// ClassroomAssigner handles assigning classrooms to lessons within the schedule generator..
type ClassroomAssigner interface {
	// Basic interface for generator component.
	GeneratorComponent[ClassroomAssignerInput, []responses.LessonWithoutClassroom]
	// Validates that classrooms can be assigned to all lessons.
	CheckAvailability(ClassroomAssignerInput) error
}

type ComponentIdentifier string

const (
	EvenWeekdayAllocatorID ComponentIdentifier = "even_weekday_allocator"

	OnePerWeekTimeSlotAssignerID ComponentIdentifier = "one_per_week_time_slot_assigner"
	BruteTimeSlotAssignerID      ComponentIdentifier = "brute_time_slot_assigner"

	MunkresClassroomAssignerID ComponentIdentifier = "munkres_classroom_assigner"
)
