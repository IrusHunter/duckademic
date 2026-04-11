package generatorResponses

type BoneLessons struct {
	Lessons []BoneLesson `json:"bone_lessons"`
	Errors  []error      `json:"errors"`
}

type BoneLesson struct {
	Teacher      CommonEntity  `json:"teacher"`
	StudentGroup CommonEntity  `json:"student_group"`
	Discipline   CommonEntity  `json:"discipline"`
	LessonType   CommonEntity  `json:"lesson_type"`
	Day          int           `json:"day"`
	Slot         int           `json:"slot"`
	Classroom    *CommonEntity `json:"classroom,omitempty"`
}
