package types

import (
	"github.com/google/uuid"
)

type Classroom struct {
	ID         uuid.UUID `json:"id"`
	RoomNumber string    `json:"room_number"`
	Capacity   int       `json:"capacity"`
}
