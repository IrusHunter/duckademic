package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/teacher_load/entities"
	"github.com/IrusHunter/duckademic/services/teacher_load/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type GroupCohortService interface {
	platform.BaseService[entities.GroupCohort]
}

func NewGroupCohortService(
	gr repositories.GroupCohortRepository,
	eb events.EventBus,
) GroupCohortService {

	sc := platform.NewServiceConfig(
		"GroupCohortService",
		filepath.Join("data", "group_cohorts.json"),
		"group_cohort",
	)

	res := &groupCohortService{
		repository: gr,
	}

	res.BaseService = platform.NewBaseService(sc, gr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.GroupCohort]{},
	)

	res.logger = res.GetLogger()

	eb.Subscribe(
		contextutil.SetTraceID(context.Background()),
		string(events.GroupCohortRT),
		res.eventHandler,
	)

	return res
}

type groupCohortService struct {
	platform.BaseService[entities.GroupCohort]
	repository repositories.GroupCohortRepository
	logger     logger.Logger
}

func (s *groupCohortService) eventHandler(ctx context.Context, b []byte) {
	gc, err := events.FromByteConvertor[events.GroupCohortRE](b)
	if err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "GroupCohortRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "GroupCohortRTHandler",
		fmt.Sprintf("received %s", gc), logger.EventDataReceived,
	)

	trueGC := entities.GroupCohort{
		ID:   gc.ID,
		Slug: gc.Slug,
		Name: gc.Name,
	}

	switch gc.Event {
	case events.EntityCreated:
		s.Add(ctx, trueGC)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, gc.ID, trueGC)
	case events.EntityDeleted:
		s.Delete(ctx, gc.ID)
	}
}

func (s *groupCohortService) Seed(ctx context.Context) error {
	groupCohorts := []entities.GroupCohort{}

	if err := jsonutil.ReadFileTo(
		filepath.Join("data", "group_cohorts.json"),
		&groupCohorts,
	); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load group cohorts seed data: %w", err),
			logger.ServiceDataFetchFailed,
		)
	}

	var lastError error

	for _, gc := range groupCohorts {
		trueGC := s.repository.FindBySlug(ctx, gc.Slug)
		if trueGC == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("group cohort with slug %q not found", gc.Slug),
				logger.ServiceDataFetchFailed,
			)
			continue
		}

		_, err := s.Update(ctx, trueGC.ID, gc)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to update %s: %w", gc, err),
				logger.ServiceRepositoryFailed,
			)
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d group cohorts updated successfully", len(groupCohorts)),
		logger.ServiceOperationSuccess,
	)

	return lastError
}

func (s *groupCohortService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	groupCohort entities.GroupCohort,
) (entities.GroupCohort, error) {

	updatedGC, err := s.repository.ExternalUpdate(ctx, id, groupCohort)
	if err != nil {
		return entities.GroupCohort{},
			s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "ExternalUpdate",
				err, logger.ServiceRepositoryFailed,
			)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedGC),
		logger.ServiceOperationSuccess,
	)

	return updatedGC, nil
}
