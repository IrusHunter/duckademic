package entities

import "github.com/google/uuid"

type Teacher struct {
	EmployeeID       uuid.UUID
	Email            string
	AcademicDegreeID uuid.UUID
	AcademicRankID   uuid.UUID
	AcademicRankStr  string
}
