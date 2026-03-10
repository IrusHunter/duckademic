package entities

import "github.com/google/uuid"

type AcademicDegree struct {
	ID   uuid.UUID
	Slug string
	Name string
}
