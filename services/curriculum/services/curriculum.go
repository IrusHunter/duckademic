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

type CurriculumService interface {
	platform.BaseService[entities.Curriculum]
}

func NewCurriculumService(
	cr repositories.CurriculumRepository,
) CurriculumService {
	sc := platform.NewServiceConfig(
		"CurriculumService",
		filepath.Join("data", "curriculums.json"),
		"curriculum",
	)

	res := &curriculumService{
		repository: cr,
	}

	res.BaseService = platform.NewBaseService(sc, cr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Curriculum]{
			platform.OnAddPrepare:   res.onAddPrepare,
			platform.ValidateEntity: res.validateEntity,
		},
	)

	res.logger = res.GetLogger()

	return res
}

type curriculumService struct {
	platform.BaseService[entities.Curriculum]
	repository repositories.CurriculumRepository
	logger     logger.Logger
}

func (s *curriculumService) validateEntity(ctx context.Context, curriculum *entities.Curriculum) error {
	if err := curriculum.ValidateName(); err != nil {
		return err
	}
	if err := curriculum.ValidateDurationYears(); err != nil {
		return err
	}
	if err := curriculum.ValidateEffectiveFrom(); err != nil {
		return err
	}

	return nil
}
func (s *curriculumService) onAddPrepare(ctx context.Context, curriculum *entities.Curriculum) error {
	slugStr := slug.Make(curriculum.Name)
	if other := s.repository.FindBySlug(ctx, slugStr); other != nil {
		return fmt.Errorf("curriculum with slug %q already exists", slugStr)
	}
	curriculum.ID = uuid.New()
	curriculum.Slug = slugStr
	return nil
}
