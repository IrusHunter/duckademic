package responses

type GeneratedLessons struct {
	Lessons []GeneratedLesson  `json:"lessons"`
	Errors  []UnassignedLesson `json:"errors"`
}

type GeneratedLesson struct {
	CommonLesson
	Days      []int         `json:"days"`
	Slot      int           `json:"slot"`
	Classroom *CommonEntity `json:"classroom,omitempty"`
}
