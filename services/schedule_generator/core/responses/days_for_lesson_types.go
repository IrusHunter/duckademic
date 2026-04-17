package responses

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
