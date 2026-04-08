package generatorResponses

type DaysForLessonTypes struct {
	StudentGroups []StudentGroupWithLTypeDays `json:"student_groups"`
	Errors        []error                     `json:"errors"`
}

type StudentGroupWithLTypeDays struct {
	CommonEntity
	WeekdayLessonTypes []LessonTypeWeekdayBinding `json:"weekday_lesson_types"`
}

type LessonTypeWeekdayBinding struct {
	CommonEntity
	Weekday int `json:"weekday"`
}
