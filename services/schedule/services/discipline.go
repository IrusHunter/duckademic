package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type DisciplineService interface {
	platform.BaseService[entities.Discipline]
}

func NewDisciplineService(dr repositories.DisciplineRepository, eb events.EventBus) DisciplineService {
	sc := platform.NewServiceConfig("DisciplineService", filepath.Join("data", "disciplines.json"), "discipline")

	res := &disciplineService{
		repository: dr,
	}
	res.BaseService = platform.NewBaseService(sc, dr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Discipline]{},
	)
	res.logger = res.GetLogger()

	eb.Subscribe(contextutil.SetTraceID(context.Background()), string(events.DisciplineRT), res.eventHandler)

	return res
}

type disciplineService struct {
	platform.BaseService[entities.Discipline]
	repository repositories.DisciplineRepository
	logger     logger.Logger
}

func (s *disciplineService) eventHandler(ctx context.Context, b []byte) {
	dr, err := events.FromByteConvertor[events.DisciplineRE](b)
	if err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "DisciplineRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "DisciplineRTHandler",
		fmt.Sprintf("received %s", dr), logger.EventDataReceived,
	)

	trueDR := entities.Discipline{
		ID:   dr.ID,
		Name: dr.Name,
		Slug: dr.Slug,
	}

	switch dr.Event {
	case events.EntityCreated:
		s.Add(ctx, trueDR)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, dr.ID, trueDR)
	case events.EntityDeleted:
		s.Delete(ctx, dr.ID)
	}
}

func (s *disciplineService) Seed(ctx context.Context) error {
	disciplines := []entities.Discipline{}
	if err := jsonutil.ReadFileTo(filepath.Join("data", "disciplines.json"), &disciplines); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load disciplines seed data: %w", err), logger.ServiceDataFetchFailed,
		)
	}

	var lastError error
	for _, discipline := range disciplines {
		tdr := s.repository.FindFirstByName(ctx, discipline.Name)
		if tdr == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("discipline with name %q not found", discipline.Name), logger.ServiceDataFetchFailed,
			)
			continue
		}

		_, err := s.Update(ctx, tdr.ID, discipline)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to update %s: %w", discipline, err), logger.ServiceRepositoryFailed,
			)
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d disciplines updated successfully", len(disciplines)), logger.ServiceOperationSuccess,
	)
	return lastError
}

func (s *disciplineService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	discipline entities.Discipline,
) (entities.Discipline, error) {
	updatedDR, err := s.repository.ExternalUpdate(ctx, id, discipline)
	if err != nil {
		return entities.Discipline{}, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "ExternalUpdate",
			err, logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedDR), logger.ServiceOperationSuccess)
	return updatedDR, nil
}
