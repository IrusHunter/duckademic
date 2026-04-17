package responses

type BoneLessons struct {
	Lessons []BoneLesson       `json:"bone_lessons"`
	Errors  []UnassignedLesson `json:"errors"`
}

type BoneLesson struct {
	CommonLesson
	Day       int           `json:"day"`
	Slot      int           `json:"slot"`
	Classroom *CommonEntity `json:"classroom,omitempty"`
}
