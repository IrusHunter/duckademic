package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/curriculum/entities"
	"github.com/IrusHunter/duckademic/services/curriculum/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
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
	eb events.EventBus,
) DisciplineService {
	sc := platform.NewServiceConfig(
		"DisciplineService",
		filepath.Join("data", "disciplines.json"),
		"discipline",
	)

	res := &disciplineService{
		repository: dr,
		eventBus:   eb,
	}

	res.BaseService = platform.NewBaseServiceWithEventBus(sc, dr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Discipline]{
			platform.OnAddPrepare:   res.onAddPrepare,
			platform.ValidateEntity: res.validateEntity,
		},
		eb,
	)

	res.logger = res.GetLogger()

	return res
}

type disciplineService struct {
	platform.BaseService[entities.Discipline]
	repository repositories.DisciplineRepository
	logger     logger.Logger
	eventBus   events.EventBus
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

func (s *disciplineService) Seed(ctx context.Context) error {
	disciplines := []entities.Discipline{}
	if err := jsonutil.ReadFileTo(filepath.Join("data", "disciplines.json"), &disciplines); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load disciplines seed data: %w", err), logger.ServiceValidationFailed,
		)
	}

	var lastError error
	for _, discipline := range disciplines {
		_, err := s.Add(ctx, discipline)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to add %s: %w", discipline.String(), err), logger.ServiceValidationFailed,
			)
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d disciplines added successfully", len(disciplines)), logger.ServiceOperationSuccess,
	)

	return lastError
}

func (s *disciplineService) Add(
	ctx context.Context, discipline entities.Discipline,
) (entities.Discipline, error) {
	addedD, err := s.BaseService.Add(ctx, discipline)
	if err == nil {
		s.sendChanges(ctx, addedD, events.EntityCreated)
	}
	return addedD, err
}
func (s *disciplineService) Delete(
	ctx context.Context, id uuid.UUID,
) (entities.Discipline, error) {
	deletedD, err := s.BaseService.Delete(ctx, id)
	if err == nil {
		s.sendChanges(ctx, deletedD, events.EntityDeleted)
	}
	return deletedD, err
}
func (s *disciplineService) Update(
	ctx context.Context, id uuid.UUID, discipline entities.Discipline,
) (entities.Discipline, error) {
	updatedD, err := s.BaseService.Update(ctx, id, discipline)
	if err == nil {
		s.sendChanges(ctx, updatedD, events.EntityUpdated)
	}
	return updatedD, err
}

func (s *disciplineService) sendChanges(
	ctx context.Context,
	discipline entities.Discipline,
	event events.EventType,
) {
	eventD := events.DisciplineRE{
		Event: event,
		ID:    discipline.ID,
		Slug:  discipline.Slug,
		Name:  discipline.Name,
	}

	s.BaseService.SendChanges(ctx, eventD, event, events.DisciplineRT)
}
