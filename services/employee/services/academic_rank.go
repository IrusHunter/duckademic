package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/employees/entities"
	"github.com/IrusHunter/duckademic/services/employees/repositories"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

// AcademicRankService provides operations to initialize and manage academic ranks.
type AcademicRankService interface {
	platform.BaseService[entities.AcademicRank]
}

// NewAcademicRankService creates a new AcademicRankService instance.
//
// It requires an academic rank repository (arr) and an event bus (eb).
func NewAcademicRankService(arr repositories.AcademicRankRepository, eb events.EventBus) AcademicRankService {
	sc := platform.NewServiceConfig("AcademicRankService", filepath.Join("data", "academic_ranks.json"), "academic_rank")

	res := &academicRankService{
		repository: arr,
	}
	res.BaseService = platform.NewBaseServiceWithEventBus(sc, arr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.AcademicRank]{
			platform.OnAddPrepare:   res.onAddPrepare,
			platform.ValidateEntity: res.validateEntity,
		}, eb,
	)
	res.logger = res.GetLogger()

	return res
}

type academicRankService struct {
	platform.BaseService[entities.AcademicRank]
	repository repositories.AcademicRankRepository
	eventBus   events.EventBus
	logger     logger.Logger
}

func (s *academicRankService) validateEntity(ctx context.Context, academicRank *entities.AcademicRank) error {
	if err := academicRank.ValidateTitle(); err != nil {
		return err
	}

	return nil
}
func (s *academicRankService) onAddPrepare(ctx context.Context, academicRank *entities.AcademicRank) error {
	slug := slug.Make(academicRank.Title)
	if other := s.repository.FindBySlug(ctx, slug); other != nil {
		return fmt.Errorf("academic rank with slug %q already exists", slug)
	}
	academicRank.ID = uuid.New()
	academicRank.Slug = slug

	return nil
}

func (s *academicRankService) Add(
	ctx context.Context, academicRank entities.AcademicRank,
) (entities.AcademicRank, error) {
	addedAR, err := s.BaseService.Add(ctx, academicRank)
	if err == nil {
		eventAR := events.AcademicRankRE{
			Event: events.EntityCreated,
			ID:    addedAR.ID,
			Title: addedAR.Title,
			Slug:  addedAR.Slug,
		}

		s.SendChanges(ctx, eventAR, events.EntityCreated, events.AcademicRankRT)
	}
	return addedAR, err
}
func (s *academicRankService) Delete(
	ctx context.Context, id uuid.UUID,
) (entities.AcademicRank, error) {
	deletedAR, err := s.BaseService.Delete(ctx, id)
	if err == nil {
		eventAR := events.AcademicRankRE{
			Event: events.EntityCreated,
			ID:    deletedAR.ID,
			Title: deletedAR.Title,
			Slug:  deletedAR.Slug,
		}

		s.SendChanges(ctx, eventAR, events.EntityDeleted, events.AcademicRankRT)
	}
	return deletedAR, err
}
func (s *academicRankService) Update(
	ctx context.Context, id uuid.UUID, academicRank entities.AcademicRank,
) (entities.AcademicRank, error) {
	updatedAR, err := s.BaseService.Update(ctx, id, academicRank)
	if err == nil {
		eventAR := events.AcademicRankRE{
			Event: events.EntityCreated,
			ID:    updatedAR.ID,
			Title: updatedAR.Title,
			Slug:  updatedAR.Slug,
		}

		s.SendChanges(ctx, eventAR, events.EntityUpdated, events.AcademicRankRT)
	}
	return updatedAR, err
}
