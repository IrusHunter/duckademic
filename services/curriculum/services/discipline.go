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

type DisciplineService interface {
	platform.BaseService[entities.Discipline]
}

func NewDisciplineService(
	dr repositories.DisciplineRepository,
) DisciplineService {
	sc := platform.NewServiceConfig(
		"DisciplineService",
		filepath.Join("data", "disciplines.json"),
		"discipline",
	)

	res := &disciplineService{
		repository: dr,
	}

	res.BaseService = platform.NewBaseService(sc, dr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Discipline]{
			platform.OnAddPrepare:   res.onAddPrepare,
			platform.ValidateEntity: res.validateEntity,
		},
	)

	res.logger = res.GetLogger()

	return res
}

type disciplineService struct {
	platform.BaseService[entities.Discipline]
	repository repositories.DisciplineRepository
	logger     logger.Logger
}

func (s *disciplineService) validateEntity(ctx context.Context, discipline *entities.Discipline) error {
	if err := discipline.ValidateName(); err != nil {
		return err
	}
	return nil
}

func (s *disciplineService) onAddPrepare(ctx context.Context, discipline *entities.Discipline) error {
	slugStr := slug.Make(discipline.Name)
	if other := s.repository.FindBySlug(ctx, slugStr); other != nil {
		return fmt.Errorf("discipline with slug %q already exists", slugStr)
	}
	discipline.ID = uuid.New()
	discipline.Slug = slugStr
	return nil
}
