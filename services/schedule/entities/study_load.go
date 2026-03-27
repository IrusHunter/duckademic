package entities

import "github.com/google/uuid"

type TeacherLoad struct {
	ID                     uuid.UUID
	TeacherID              uuid.UUID
	LessonTypeAssignmentID uuid.UUID
	StudentGroupID         uuid.UUID
}
