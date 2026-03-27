package entities

import "github.com/google/uuid"

type TeacherLoad struct {
	ID            uuid.UUID
	TeacherID     uuid.UUID
	DisciplineID  uuid.UUID
	LessonTypeID  uuid.UUID
	GroupCohortID uuid.UUID
	GroupCount    int
}
