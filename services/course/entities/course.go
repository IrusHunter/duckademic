package entities

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID           uuid.UUID
	DisciplineID uuid.UUID
	Manager      uuid.UUID
	Name         string
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
