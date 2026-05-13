package responses

import "github.com/google/uuid"

type CommonEntity struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func FormCommonEntity(id uuid.UUID, name string) CommonEntity {
	return CommonEntity{
		ID:   id,
		Name: name,
	}
}

type UnassignedLesson struct {
	StudyLoad
	Count int
}
