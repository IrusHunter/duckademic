package services

import (
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	externalEntities "github.com/IrusHunter/duckademic/services/schedule_generator/entities"
	"github.com/google/uuid"
)

// LessonTypeService aggregates and manages lesson types that the generator works with.
type LessonTypeService interface {
	Find(uuid.UUID) *entities.LessonType // Returns a pointer to the lesson type with the given ID.
	GetAll() []*entities.LessonType      // Returns an array with all lesson types as pointers.
}

// NewLessonTypeService creates a new basic LessonTypeService instance.
//
// It requires an array of database lesson types (lt).
func NewLessonTypeService(lt []externalEntities.LessonType) (LessonTypeService, error) {
	lts := lessonTypeService{
		lessonTypes: make([]*entities.LessonType, len(lt)),
	}

	for i, lt := range lt {
		lts.lessonTypes[i] = &entities.LessonType{
			ID:    lt.ID,
			Name:  lt.Name,
			Weeks: lt.ReservedWeeks,
			Value: lt.HoursValue,
		}
	}

	return &lts, nil
}

// lessonTypeService is the basic implementation of the LessonTypeService interface.
type lessonTypeService struct {
	lessonTypes []*entities.LessonType
}

func (lts *lessonTypeService) Find(id uuid.UUID) *entities.LessonType {
	for i := range lts.lessonTypes {
		if lts.lessonTypes[i].ID == id {
			return lts.lessonTypes[i]
		}
	}
	return nil
}
func (lts *lessonTypeService) GetAll() []*entities.LessonType {
	return lts.lessonTypes
}
