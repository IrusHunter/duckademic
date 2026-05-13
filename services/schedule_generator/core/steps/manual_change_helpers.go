package steps

import (
	"fmt"
	"strconv"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/services"
	"github.com/google/uuid"
)

func getUUID(data map[string]string, key string) (uuid.UUID, error) {
	v, ok := data[key]
	if !ok {
		return uuid.UUID{}, fmt.Errorf("data doesn't contain key %q", key)
	}

	id, err := uuid.Parse(v)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to parse %q as uuid.UUID: %w", key, err)
	}

	return id, nil
}
func getInt(data map[string]string, key string) (int, error) {
	v, ok := data[key]
	if !ok {
		return 0, fmt.Errorf("data doesn't contain key %q", key)
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %q as int: %w", key, err)
	}

	return i, nil
}

func FormWeekdayBindingOverride(data map[string]string) (WeekdayBindingOverride, error) {
	var res WeekdayBindingOverride
	var err error

	res.StudentGroupID, err = getUUID(data, "student_group_id")
	if err != nil {
		return res, err
	}

	res.LessonTypeID, err = getUUID(data, "lesson_type_id")
	if err != nil {
		return res, err
	}

	if v, err := getInt(data, "old_weekday"); err == nil {
		res.OldWeekday = &v
	}

	if v, err := getInt(data, "new_weekday"); err == nil {
		res.NewWeekday = &v
	}

	return res, nil
}
func FormTimeSlotForLessonOverride(data map[string]string) (TimeSlotForLessonOverride, error) {
	var res TimeSlotForLessonOverride
	var err error

	res.StudyLoadID, err = getUUID(data, "study_load_id")
	if err != nil {
		return res, err
	}

	if v, err := getInt(data, "old_day"); err == nil {
		res.OldDay = &v
	}

	if v, err := getInt(data, "old_slot"); err == nil {
		res.OldSlot = &v
	}

	if v, err := getInt(data, "new_day"); err == nil {
		res.NewDay = &v
	}

	if v, err := getInt(data, "new_slot"); err == nil {
		res.NewSlot = &v
	}

	return res, nil
}
func FormClassroomForLessonOverride(
	data map[string]string,
) (ClassroomForLessonOverride, error) {
	var res ClassroomForLessonOverride
	var err error

	res.LessonID, err = getUUID(data, "lesson_id")
	if err != nil {
		return res, err
	}

	if v, err := getUUID(data, "new_classroom_id"); err == nil {
		res.NewClassroomID = &v
	}

	return res, nil
}

func TimeSlotForLessonOverrideChange(
	data map[string]string,
	lessonService services.LessonService,
	studyLoadService services.StudyLoadService,
) error {
	timeSlotForLessonOverride, err := FormTimeSlotForLessonOverride(data)
	if err != nil {
		return fmt.Errorf("failed to parse time slot for lesson override: %w", err)
	}

	studyLoad := studyLoadService.FindByID(timeSlotForLessonOverride.StudyLoadID)
	if studyLoad == nil {
		return fmt.Errorf("study load with id %q not found", timeSlotForLessonOverride.StudyLoadID)
	}

	if timeSlotForLessonOverride.OldDay != nil && timeSlotForLessonOverride.OldSlot != nil {
		oldDay := *timeSlotForLessonOverride.OldDay
		oldSlot := *timeSlotForLessonOverride.OldSlot
		oldLessonSlot := entities.NewLessonSlot(oldDay, oldSlot)

		lesson := lessonService.FindByStudyLoadIDAndSlot(studyLoad.ID, oldLessonSlot)
		if lesson == nil {
			return fmt.Errorf(
				"lesson for study load with id %q and at %s not found",
				timeSlotForLessonOverride.StudyLoadID,
				oldLessonSlot.String(),
			)
		}

		// move lesson if all data was given
		if timeSlotForLessonOverride.NewDay != nil && timeSlotForLessonOverride.NewSlot != nil {
			newDay := *timeSlotForLessonOverride.NewDay
			newSlot := *timeSlotForLessonOverride.NewSlot
			newLessonSlot := entities.NewLessonSlot(newDay, newSlot)

			if err := lessonService.MoveLessonTo(lesson, newLessonSlot); err != nil {
				return err
			}
		} else { // remove slot for lesson
			lessonService.UnassignLesson(lesson)
		}
	} else {
		// assign lesson slot for study load
		if timeSlotForLessonOverride.NewDay != nil && timeSlotForLessonOverride.NewSlot != nil {
			newDay := *timeSlotForLessonOverride.NewDay
			newSlot := *timeSlotForLessonOverride.NewSlot
			newLessonSlot := entities.NewLessonSlot(newDay, newSlot)

			err := lessonService.AssignLesson(studyLoad, newLessonSlot)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
func ClassroomForLessonOverrideChange(
	data map[string]string,
	lessonService services.LessonService,
	classroomService services.ClassroomService,
) error {
	classroomForLessonOverride, err := FormClassroomForLessonOverride(data)
	if err != nil {
		return fmt.Errorf("failed to parse classroom for lesson override: %w", err)
	}

	lesson := lessonService.FindByID(classroomForLessonOverride.LessonID)
	if lesson == nil {
		return fmt.Errorf("lesson with id %q not found", classroomForLessonOverride.LessonID)
	}

	if classroomForLessonOverride.NewClassroomID != nil {
		newClassroomID := *classroomForLessonOverride.NewClassroomID

		classroom := classroomService.Find(newClassroomID)
		if classroom == nil {
			return fmt.Errorf("classroom with id %q not found", newClassroomID)
		}

		err := lesson.SetClassroom(classroom)
		if err != nil {
			return err
		}
	} else {
		err := lesson.RemoveClassroom()
		if err != nil {
			return err
		}
	}

	return nil
}
