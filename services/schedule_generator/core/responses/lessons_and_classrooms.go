package responses

type GeneratedLessonsWithC struct {
	LessonsWithClassroom    []GeneratedLesson        `json:"lessons_with_classroom"`
	LessonsWithoutClassroom []LessonWithoutClassroom `json:"lessons_without_classroom"`
}

type LessonWithoutClassroom struct {
	CommonLesson
	Day  int `json:"day"`
	Slot int `json:"slot"`
}
