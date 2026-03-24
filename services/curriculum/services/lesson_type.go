package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/curriculum/entities"
	"github.com/IrusHunter/duckademic/services/curriculum/repositories"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type LessonTypeService interface {
	platform.BaseService[entities.LessonType]
}

func NewLessonTypeService(
	ltr repositories.LessonTypeRepository,
) LessonTypeService {
	sc := platform.NewServiceConfig(
		"LessonTypeService",
		filepath.Join("data", "lesson_types.json"),
		"lesson_type",
	)

	res := &lessonTypeService{
		repository: ltr,
	}

	res.BaseService = platform.NewBaseService(sc, ltr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.LessonType]{
			platform.OnAddPrepare:   res.onAddPrepare,
			platform.ValidateEntity: res.validateEntity,
		},
	)

	res.logger = res.GetLogger()

	return res
}

type lessonTypeService struct {
	platform.BaseService[entities.LessonType]
	repository repositories.LessonTypeRepository
	logger     logger.Logger
}

func (s *lessonTypeService) validateEntity(ctx context.Context, lt *entities.LessonType) error {
	if err := lt.ValidateName(); err != nil {
		return err
	}
	if err := lt.ValidateHoursValue(); err != nil {
		return err
	}

	return nil
}
func (s *lessonTypeService) onAddPrepare(ctx context.Context, lt *entities.LessonType) error {
	slugStr := slug.Make(lt.Name)
	if other := s.repository.FindBySlug(ctx, slugStr); other != nil {
		return fmt.Errorf("lesson type with slug %q already exists", slugStr)
	}
	lt.ID = uuid.New()
	lt.Slug = slugStr
	return nil
}
