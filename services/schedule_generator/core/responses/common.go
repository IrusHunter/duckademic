package responses

import "github.com/google/uuid"

type CommonEntity struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type CommonLesson struct {
	Teacher      CommonEntity `json:"teacher"`
	StudentGroup CommonEntity `json:"student_group"`
	Discipline   CommonEntity `json:"discipline"`
	LessonType   CommonEntity `json:"lesson_type"`
}

type UnassignedLesson struct {
	CommonLesson
	Count int
}
