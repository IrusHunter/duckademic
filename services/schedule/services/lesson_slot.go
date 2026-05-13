package services

import (
	"context"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/repositories"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type LessonSlotService interface {
	platform.BaseService[entities.LessonSlot]
}

func NewLessonSlotService(
	lr repositories.LessonSlotRepository,
) LessonSlotService {
	sc := platform.NewServiceConfig(
		"LessonSlotService",
		filepath.Join("data", "lesson_slots.json"),
		"lesson slot",
	)

	res := &lessonSlotService{
		repository: lr,
	}

	res.BaseService = platform.NewBaseService(
		sc,
		lr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.LessonSlot]{
			platform.OnAddPrepare:       res.onAddPrepare,
			platform.ValidateEntity:     res.validateEntity,
			platform.OnUpdateValidation: res.validateEntity,
		},
	)

	return res
}

type lessonSlotService struct {
	platform.BaseService[entities.LessonSlot]
	repository repositories.LessonSlotRepository
}

func (s *lessonSlotService) validateEntity(ctx context.Context, ls *entities.LessonSlot) error {
	if err := ls.ValidateDuration(); err != nil {
		return err
	}
	if err := ls.ValidateStartTime(); err != nil {
		return err
	}

	return nil
}

func (s *lessonSlotService) onAddPrepare(ctx context.Context, ls *entities.LessonSlot) error {
	ls.ID = uuid.New()

	return nil
}
