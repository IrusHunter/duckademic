package responses

import "github.com/google/uuid"

type StudyLoad struct {
	ID             uuid.UUID `json:"id"`
	TeacherID      uuid.UUID `json:"teacher_id"`
	StudentGroupID uuid.UUID `json:"student_group_id"`
	DisciplineID   uuid.UUID `json:"discipline_id"`
	LessonTypeID   uuid.UUID `json:"lesson_type_id"`
}
