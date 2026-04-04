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

type StudentGroupService interface {
	platform.BaseService[entities.StudentGroup]
}

func NewStudentGroupService(gr repositories.StudentGroupRepository, eb events.EventBus) StudentGroupService {
	sc := platform.NewServiceConfig("StudentGroupService", filepath.Join("data", "student_groups.json"), "student group")

	res := &studentGroupService{
		repository: gr,
	}
	res.BaseService = platform.NewBaseService(sc, gr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.StudentGroup]{},
	)
	res.logger = res.GetLogger()

	eb.Subscribe(contextutil.SetTraceID(context.Background()), string(events.StudentGroupRT), res.eventHandler)

	return res
}

type studentGroupService struct {
	platform.BaseService[entities.StudentGroup]
	repository repositories.StudentGroupRepository
	logger     logger.Logger
}

func (s *studentGroupService) eventHandler(ctx context.Context, b []byte) {
	groupEvent, err := events.FromByteConvertor[events.StudentGroupRE](b)
	if err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "StudentGroupRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "StudentGroupRTHandler",
		fmt.Sprintf("received %s", groupEvent), logger.EventDataReceived,
	)

	trueG := entities.StudentGroup{
		ID:            groupEvent.ID,
		Slug:          groupEvent.Slug,
		Name:          groupEvent.Name,
		GroupCohortID: groupEvent.GroupCohortID,
	}

	switch groupEvent.Event {
	case events.EntityCreated:
		s.Add(ctx, trueG)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, groupEvent.ID, trueG)
	case events.EntityDeleted:
		s.Delete(ctx, groupEvent.ID)
	}
}

func (s *studentGroupService) Seed(ctx context.Context) error {
	groups := []entities.StudentGroup{}
	if err := jsonutil.ReadFileTo(filepath.Join("data", "student_groups.json"), &groups); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load student groups seed data: %w", err), logger.ServiceDataFetchFailed,
		)
	}

	var lastError error
	for _, group := range groups {
		trueG := s.repository.FindFirstByName(ctx, group.Name)
		if trueG == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("student group with name %q not found", group.Name), logger.ServiceDataFetchFailed,
			)
			continue
		}

		_, err := s.Update(ctx, trueG.ID, group)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to update %s: %w", group, err), logger.ServiceRepositoryFailed,
			)
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d student groups updated successfully", len(groups)), logger.ServiceOperationSuccess,
	)
	return lastError
}

func (s *studentGroupService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	group entities.StudentGroup,
) (entities.StudentGroup, error) {
	updatedG, err := s.repository.ExternalUpdate(ctx, id, group)
	if err != nil {
		return entities.StudentGroup{}, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "ExternalUpdate",
			err, logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedG), logger.ServiceOperationSuccess)
	return updatedG, nil
}
