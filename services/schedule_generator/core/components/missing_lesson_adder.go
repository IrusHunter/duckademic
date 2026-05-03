package components

import (
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
)

// NewMissingLessonAdder creates a MissingLessonsAdder instance.
func NewMissingLessonAdder() TimeSlotAssigner {
	return &missingLessonsAdder{errorService: NewErrorService[responses.UnassignedLesson, *MissingLessonsAdderError]()}
}

type missingLessonsAdder struct {
	errorService ErrorService[responses.UnassignedLesson, *MissingLessonsAdderError]
}

func (ma *missingLessonsAdder) Run(input TimeSlotAssignerInput) []responses.UnassignedLesson {
	for _, load := range input.StudyLoads {
		teacher := load.Teacher
		studentGroup := load.StudentGroup
		lessonType := load.Type
		// discipline := load.Discipline

		currentDay := 0
		outOfGrid := false
		for !outOfGrid {
			err := studentGroup.CheckDay(currentDay)
			if err != nil {
				outOfGrid = true
				//continue
				break
			}

			for i := range teacher.BusyGrid.Grid[currentDay] {
				slot := entities.LessonSlot{
					Day:  currentDay,
					Slot: i,
				}
				err := input.LessonService.AssignLesson(load, slot)
				if err == nil {
					currentDay -= currentDay%7 - 6
					break
				}
			}
			delta := studentGroup.GetNextDayOfType(lessonType, currentDay+1)
			if delta == -1 {
				outOfGrid = true
				continue
			}
			currentDay += delta
		}

		if !load.IsEnoughHours() {
			ma.errorService.AddError(&MissingLessonsAdderError{
				StudyLoad: load,
				Count:     load.GetRequiredSlots(),
			})
		}
	}

	return ma.errorService.GetGeneratorResponseErrors()
}
func (ma *missingLessonsAdder) GetComponentIdentifier() ComponentIdentifier {
	return BruteTimeSlotAssignerID
}

// MissingLessonsAdderError indicates that the MissingLessonsAdder failed to
// find free slot in the grids for missing lesson.
type MissingLessonsAdderError struct {
	*entities.StudyLoad
	Count int
}

func (e *MissingLessonsAdderError) Error() string {
	return fmt.Sprintf("Not enough space of %s or %s for %s %s.",
		e.StudentGroup.Name, e.Teacher.UserName, e.Type.Name, e.Discipline.Name)
}
func (e *MissingLessonsAdderError) GeneratorResponseError() responses.UnassignedLesson {
	return responses.UnassignedLesson{
		CommonLesson: responses.CommonLesson{
			Teacher: responses.CommonEntity{
				ID:   e.Teacher.ID,
				Name: e.Teacher.UserName,
			},
			StudentGroup: responses.CommonEntity{
				ID:   e.StudentGroup.ID,
				Name: e.StudentGroup.Name,
			},
			Discipline: responses.CommonEntity{
				ID:   e.Discipline.ID,
				Name: e.Discipline.Name,
			},
			LessonType: responses.CommonEntity{
				ID:   e.Type.ID,
				Name: e.Type.Name,
			},
		},
		Count: e.Count,
	}
}
