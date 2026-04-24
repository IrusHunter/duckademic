package entities

import (
	"time"

	"github.com/google/uuid"
)

type StudentCourse struct {
	ID        uuid.UUID
	CourseID  uuid.UUID
	StudentID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
