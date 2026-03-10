package entities

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	ID          uuid.UUID
	Slug        string
	FirstName   string
	LastName    string
	MiddleName  string
	AvatarURL   string
	PhoneNumber string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}
