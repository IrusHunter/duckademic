package generatorResponses

import "github.com/google/uuid"

type CommonEntity struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
