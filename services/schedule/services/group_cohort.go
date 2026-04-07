package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type GroupCohortService interface {
	platform.BaseService[entities.GroupCohort]
}

func NewGroupCohortService(gr repositories.GroupCohortRepository, eb events.EventBus) GroupCohortService {
	sc := platform.NewServiceConfig(
		"GroupCohortService",
		filepath.Join("data", "group_cohorts.json"),
		"group cohort",
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
	cohortEvent, err := events.FromByteConvertor[events.GroupCohortRE](b)
	if err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "GroupCohortRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "GroupCohortRTHandler",
		fmt.Sprintf("received %s", cohortEvent), logger.EventDataReceived,
	)

	trueC := entities.GroupCohort{
		ID:   cohortEvent.ID,
		Slug: cohortEvent.Slug,
		Name: cohortEvent.Name,
	}

	switch cohortEvent.Event {
	case events.EntityCreated:
		s.Add(ctx, trueC)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, cohortEvent.ID, trueC)
	case events.EntityDeleted:
		s.Delete(ctx, cohortEvent.ID)
	}
}

func (s *groupCohortService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	cohort entities.GroupCohort,
) (entities.GroupCohort, error) {
	updatedC, err := s.repository.ExternalUpdate(ctx, id, cohort)
	if err != nil {
		return entities.GroupCohort{}, s.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"ExternalUpdate",
			err,
			logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedC),
		logger.ServiceOperationSuccess,
	)

	return updatedC, nil
}
