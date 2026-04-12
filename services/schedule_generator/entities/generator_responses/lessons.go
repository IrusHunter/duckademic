package generatorResponses

type Lessons struct {
	Lessons []Lesson `json:"lessons"`
	Errors  []error  `json:"errors"`
}

type Lesson struct {
	Teacher      CommonEntity  `json:"teacher"`
	StudentGroup CommonEntity  `json:"student_group"`
	Discipline   CommonEntity  `json:"discipline"`
	LessonType   CommonEntity  `json:"lesson_type"`
	Days         []int         `json:"days"`
	Slot         int           `json:"slot"`
	Classroom    *CommonEntity `json:"classroom,omitempty"`
}
