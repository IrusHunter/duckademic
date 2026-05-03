package responses

import (
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
)

type DaysForLessonTypes struct {
	StudentGroups []StudentGroupWithLTypeDays `json:"student_groups"`
	Errors        []LessonTypeDayDebt         `json:"errors"`
}

type StudentGroupWithLTypeDays struct {
	CommonEntity
	WeekdayLessonTypes []LessonTypeWeekdayBinding `json:"weekday_lesson_types"`
}

type LessonTypeWeekdayBinding struct {
	CommonEntity
	Weekday int `json:"weekday"`
}

type LessonTypeDayDebt struct {
	StudentGroup CommonEntity `json:"student_group"`
	LessonType   CommonEntity `json:"lesson_type"`
	SlotsDept    float64      `json:"slots_dept"`
}

func FormDaysForLessonTypes(studentGroups []*entities.StudentGroup) []StudentGroupWithLTypeDays {
	result := make([]StudentGroupWithLTypeDays, 0, len(studentGroups))

	for _, studentGroup := range studentGroups {
		resultSG := StudentGroupWithLTypeDays{
			CommonEntity: CommonEntity{
				ID:   studentGroup.ID,
				Name: studentGroup.Name,
			},
			WeekdayLessonTypes: make([]LessonTypeWeekdayBinding, 0),
		}

		for j := range 7 {
			lessonType := studentGroup.GetTypeOfDay(j)
			if lessonType != nil {
				resultSG.WeekdayLessonTypes = append(resultSG.WeekdayLessonTypes,
					LessonTypeWeekdayBinding{
						CommonEntity: CommonEntity{
							ID:   lessonType.ID,
							Name: lessonType.Name,
						},
						Weekday: j,
					})
			}
		}

		result = append(result, resultSG)
	}

	return result
}
