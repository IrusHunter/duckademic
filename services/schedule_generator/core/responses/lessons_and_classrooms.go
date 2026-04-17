package responses

type BoneLessonsWithC struct {
	LessonsWithClassroom    []BoneLesson             `json:"lessons_with_classroom"`
	LessonsWithoutClassroom []LessonWithoutClassroom `json:"lessons_without_classroom"`
}

type GeneratedLessonsWithC struct {
	LessonsWithClassroom    []GeneratedLesson        `json:"lessons_with_classroom"`
	LessonsWithoutClassroom []LessonWithoutClassroom `json:"lessons_without_classroom"`
}

type LessonWithoutClassroom struct {
	CommonLesson
	Day  int `json:"day"`
	Slot int `json:"slot"`
}
