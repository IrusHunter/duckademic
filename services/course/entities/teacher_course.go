package entities

import (
	"time"

	"github.com/google/uuid"
)

type TeacherCourse struct {
	ID        uuid.UUID
	CourseID  uuid.UUID
	TeacherID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
