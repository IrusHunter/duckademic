package responses

import "github.com/google/uuid"

type StudyLoad struct {
	ID               uuid.UUID `json:"id"`
	TeacherID        uuid.UUID `json:"teacher_id"`
	TeacherName      string    `json:"teacher_name"`
	StudentGroupID   uuid.UUID `json:"student_group_id"`
	StudentGroupName string    `json:"student_group_name"`
	DisciplineID     uuid.UUID `json:"discipline_id"`
	DisciplineName   string    `json:"discipline_name"`
	LessonTypeID     uuid.UUID `json:"lesson_type_id"`
	LessonTypeName   string    `json:"lesson_type_name"`
}
