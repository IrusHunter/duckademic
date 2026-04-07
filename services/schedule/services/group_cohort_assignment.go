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

type GroupCohortAssignmentService interface {
	platform.BaseService[entities.GroupCohortAssignment]
}

func NewGroupCohortAssignmentService(
	gr repositories.GroupCohortAssignmentRepository,
	eb events.EventBus,
) GroupCohortAssignmentService {
	sc := platform.NewServiceConfig(
		"GroupCohortAssignmentService",
		filepath.Join("data", "group_cohort_assignments.json"),
		"group cohort assignment",
	)

	res := &groupCohortAssignmentService{
		repository: gr,
	}
	res.BaseService = platform.NewBaseService(sc, gr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.GroupCohortAssignment]{},
	)
	res.logger = res.GetLogger()

	eb.Subscribe(
		contextutil.SetTraceID(context.Background()),
		string(events.GroupCohortAssignmentRT),
		res.eventHandler,
	)

	return res
}

type groupCohortAssignmentService struct {
	platform.BaseService[entities.GroupCohortAssignment]
	repository repositories.GroupCohortAssignmentRepository
	logger     logger.Logger
}

func (s *groupCohortAssignmentService) eventHandler(ctx context.Context, b []byte) {
	assignEvent, err := events.FromByteConvertor[events.GroupCohortAssignmentRE](b)
	if err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "GroupCohortAssignmentRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "GroupCohortAssignmentRTHandler",
		fmt.Sprintf("received %s", assignEvent), logger.EventDataReceived,
	)

	trueA := entities.GroupCohortAssignment{
		ID:            assignEvent.ID,
		GroupCohortID: assignEvent.GroupCohortID,
		DisciplineID:  assignEvent.DisciplineID,
		LessonTypeID:  assignEvent.LessonTypeID,
	}

	switch assignEvent.Event {
	case events.EntityCreated:
		s.Add(ctx, trueA)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, assignEvent.ID, trueA)
	case events.EntityDeleted:
		s.Delete(ctx, assignEvent.ID)
	}
}

func (s *groupCohortAssignmentService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	assignment entities.GroupCohortAssignment,
) (entities.GroupCohortAssignment, error) {
	updatedA, err := s.repository.ExternalUpdate(ctx, id, assignment)
	if err != nil {
		return entities.GroupCohortAssignment{}, s.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"ExternalUpdate",
			err,
			logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedA),
		logger.ServiceOperationSuccess,
	)

	return updatedA, nil
}
