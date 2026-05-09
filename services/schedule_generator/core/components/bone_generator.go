package components

import (
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
)

func NewBoneGenerator() TimeSlotAssigner {
	return &boneGenerator{}
}

type boneGenerator struct{}

// GenerateBoneLessons allocates lesson slots for the bone week.
// Uses brute force method, starts with teachers, then discipline and student groups,
// then free slots for lesson type.
func (bg *boneGenerator) Run(input TimeSlotAssignerInput) []responses.UnassignedLesson {
	errorService := NewErrorService[responses.UnassignedLesson, *BoneWeekError]()

	for _, load := range input.StudyLoads {
		teacher := load.Teacher
		studentGroup := load.StudentGroup
		lessonType := load.Type
		// discipline := load.Discipline

		offset := 0
		success := false

		for !success {
			// отримуємо доступний лекційний день
			day := studentGroup.GetNextDayOfType(lessonType, offset)
			if day > 7 || day < 0 {
				// якщо день був не на кістковому тижні, виникає виняток, який треба обробити якось
				errorService.AddError(&BoneWeekError{StudyLoad: load})
				break
			}

			// отримання вільного слота для групи та викладача
			lessonSlot := teacher.GetOptimalFreeSlot(studentGroup.GetFreeSlots(day), day)

			if lessonSlot != -1 {
				slot := entities.LessonSlot{Day: day, Slot: lessonSlot}
				err := input.LessonService.AssignLesson(load, slot)
				if err != nil {
					NewUnexpectedError("slot is busy but algorithm determined it as free",
						"boneGenerator", "GenerateBoneLessons", &FalseFreeSlotError{
							UnassignedLesson: load.UnassignedLesson,
							slot:             slot,
							err:              err,
						})
				}
				success = true
			}
			offset = day + 1
		}

	}
	return errorService.GetGeneratorResponseErrors()
}
func (bg *boneGenerator) GetComponentIdentifier() ComponentIdentifier {
	return OnePerWeekTimeSlotAssignerID
}

// ==========================================================================================================
// ================================================= ERRORS =================================================
// ==========================================================================================================

// BoneWeekError indicates that the BoneGenerator failed to allocate
// enough space for lessons within the bone week.
type BoneWeekError struct {
	*entities.StudyLoad
}

func (e *BoneWeekError) Error() string {
	return fmt.Sprintf("Not enough space in bone week of %s or %s for %s %s.",
		e.StudentGroup.Name, e.Teacher.UserName, e.Type.Name, e.Discipline.Name)
}
func (e *BoneWeekError) GeneratorResponseError() responses.UnassignedLesson {
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
		Count: 1,
	}
}

// FalseFreeSlotError indicates that slot is busy but algorithm determined it as free.
type FalseFreeSlotError struct {
	entities.UnassignedLesson
	slot entities.LessonSlot
	err  error
}

func (e *FalseFreeSlotError) Error() string {
	return fmt.Sprintf("false free slot %d/%d of %s or %s grid for %s %s. error: %s", e.slot.Day, e.slot.Slot,
		e.StudentGroup.Name, e.Teacher.UserName, e.Type.Name, e.Discipline.Name, e.err.Error())
}
