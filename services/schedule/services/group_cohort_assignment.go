package services

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

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
	GetByGroupCohortIDs(context.Context, []uuid.UUID) ([]entities.GroupCohortAssignment, error)
	ToGeneratorGroupCohortAssignments(context.Context, []entities.GroupCohortAssignment) []GeneratorGroupCohortAssignment
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

func (s *groupCohortAssignmentService) GetByGroupCohortIDs(
	ctx context.Context, cohortIDs []uuid.UUID,
) ([]entities.GroupCohortAssignment, error) {
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup
	var mu sync.Mutex

	var result []entities.GroupCohortAssignment
	var lastError error

	for _, cohortID := range cohortIDs {
		wg.Add(1)
		sem <- struct{}{}

		go func(cohortID uuid.UUID) {
			defer wg.Done()
			defer func() { <-sem }()

			groupCohortAssignments, err := s.repository.GetByGroupCohortID(ctx, cohortID)
			if err != nil {
				mu.Lock()
				lastError = s.GetLogger().LogAndReturnError(
					contextutil.GetTraceID(ctx),
					"GetByGroupCohortIDs",
					err,
					logger.ServiceRepositoryFailed,
				)
				mu.Unlock()
				return
			}

			mu.Lock()
			result = append(result, groupCohortAssignments...)
			mu.Unlock()

		}(cohortID)
	}

	wg.Wait()

	return result, lastError
}

type GeneratorGroupCohortAssignment struct {
	ID            uuid.UUID `json:"id"`
	GroupCohortID uuid.UUID `json:"group_cohort_id"`
	DisciplineID  uuid.UUID `json:"discipline_id"`
	LessonTypeID  uuid.UUID `json:"lesson_type_id"`
}

func (s *groupCohortAssignmentService) ToGeneratorGroupCohortAssignments(
	ctx context.Context, gc []entities.GroupCohortAssignment,
) []GeneratorGroupCohortAssignment {
	res := make([]GeneratorGroupCohortAssignment, 0, len(gc))

	for _, groupCohortAssignment := range gc {
		res = append(res, GeneratorGroupCohortAssignment{
			ID:            groupCohortAssignment.ID,
			GroupCohortID: groupCohortAssignment.GroupCohortID,
			DisciplineID:  groupCohortAssignment.DisciplineID,
			LessonTypeID:  groupCohortAssignment.LessonTypeID,
		})
	}

	return res
}
