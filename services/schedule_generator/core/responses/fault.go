package responses

import (
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
)

type Fault struct {
	TotalValue float64 `json:"total_value"`

	TeacherWindows                FaultParam[TeacherWindow]              `json:"teacher_windows"`
	StudentGroupWindows           FaultParam[StudentGroupWindow]         `json:"student_group_windows"`
	TeacherLessonOverlapping      FaultParam[TeacherLessonOverlap]       `json:"teacher_lesson_overlapping"`
	StudentGroupLessonOverlapping FaultParam[StudentGroupLessonOverlap]  `json:"student_group_lesson_overlapping"`
	ClassroomLessonOverlapping    FaultParam[ClassroomLessonOverlap]     `json:"classroom_lesson_overlapping"`
	StudentGroupOvertimeLessons   FaultParam[StudentGroupOvertimeLesson] `json:"student_group_overtime_lessons"`
	StudentGroupInvalidLessons    FaultParam[StudentGroupInvalidLesson]  `json:"student_group_invalid_lessons_by_type"`
	LessonsWithoutClassroom       FaultParam[GeneratedLesson]            `json:"lessons_without_classroom"`
	ClassroomWithOverflow         FaultParam[ClassroomOverflow]          `json:"classroom_with_overflow"`
	StudyLoadHoursDeficit         FaultParam[StudyLoadHoursDeficit]      `json:"study_load_hours_deficit"`
}

type FaultParam[T any] struct {
	Value      float64 `json:"value"`
	Defections []T     `json:"defections"`
}

type TeacherWindow struct {
	CommonEntity
	Windows []TimeSlot `json:"windows"`
}
type StudentGroupWindow struct {
	CommonEntity
	Windows []TimeSlot `json:"windows"`
}

type TeacherLessonOverlap struct {
	CommonEntity
	OverlapLessons []TimeSlot `json:"overlap_lessons"`
}
type StudentGroupLessonOverlap struct {
	CommonEntity
	OverlapLessons []TimeSlot `json:"overlap_lessons"`
}
type ClassroomLessonOverlap struct {
	CommonEntity
	OverlapLessons []TimeSlot `json:"overlap_lessons"`
}

type StudentGroupOvertimeLesson struct {
	CommonEntity
	OvertimeDays []int `json:"overtime_days"`
}

type StudentGroupInvalidLesson struct {
	CommonEntity
	InvalidLessons []InvalidLessonByType `json:"invalid_lessons"`
}

type InvalidLessonByType struct {
	TimeSlot
	LessonType       CommonEntity `json:"lesson_type"`
	ActualLessonType CommonEntity `json:"actual_lesson_type"`
}

func FormInvalidLessonByType(ils []entities.InvalidLessonByType) (res []InvalidLessonByType) {
	for _, il := range ils {
		res = append(res, InvalidLessonByType{
			TimeSlot: TimeSlot{
				Slot: il.Lesson.Slot,
				Day:  il.Lesson.Day,
			},
			LessonType:       FormCommonEntity(il.Lesson.Type.ID, il.Lesson.Type.Name),
			ActualLessonType: FormCommonEntity(il.ActualLessonType.ID, il.ActualLessonType.Name),
		})
	}
	return
}

type StudyLoadHoursDeficit struct {
	StudyLoad
	Missing int `json:"missing_hours"`
}

type ClassroomOverflow struct {
	CommonEntity
	Overflows []ClassroomOverflowItem `json:"overflows"`
	Capacity  int                     `json:"capacity"`
}

type ClassroomOverflowItem struct {
	TimeSlot
	Required int `json:"required"`
}

func FormClassroomOverflowItems(lessons []*entities.Lesson) []ClassroomOverflowItem {
	res := make([]ClassroomOverflowItem, 0, len(lessons))

	for _, lesson := range lessons {
		res = append(res, ClassroomOverflowItem{
			TimeSlot: TimeSlot{
				Day:  lesson.Day,
				Slot: lesson.Slot,
			},
			Required: lesson.StudentGroup.StudentNumber,
		})
	}

	return res
}

type TimeSlot struct {
	Day  int `json:"day"`
	Slot int `json:"slot"`
}

func FormTimeSlots(slots []entities.LessonSlot) []TimeSlot {
	res := make([]TimeSlot, 0, len(slots))
	for _, slot := range slots {
		res = append(res, TimeSlot{
			Slot: slot.Slot,
			Day:  slot.Day,
		})
	}
	return res
}
