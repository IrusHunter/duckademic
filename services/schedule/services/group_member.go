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

type GroupMemberService interface {
	platform.BaseService[entities.GroupMember]
}

func NewGroupMemberService(gr repositories.GroupMemberRepository, eb events.EventBus) GroupMemberService {
	sc := platform.NewServiceConfig(
		"GroupMemberService",
		filepath.Join("data", "group_members.json"),
		"group member",
	)

	res := &groupMemberService{
		repository: gr,
	}
	res.BaseService = platform.NewBaseService(sc, gr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.GroupMember]{},
	)
	res.logger = res.GetLogger()

	eb.Subscribe(
		contextutil.SetTraceID(context.Background()),
		string(events.GroupMemberRT),
		res.eventHandler,
	)

	return res
}

type groupMemberService struct {
	platform.BaseService[entities.GroupMember]
	repository repositories.GroupMemberRepository
	logger     logger.Logger
}

func (s *groupMemberService) eventHandler(ctx context.Context, b []byte) {
	memberEvent, err := events.FromByteConvertor[events.GroupMemberRE](b)
	if err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "GroupMemberRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "GroupMemberRTHandler",
		fmt.Sprintf("received %s", memberEvent), logger.EventDataReceived,
	)

	trueM := entities.GroupMember{
		ID:           memberEvent.ID,
		StudentID:    memberEvent.StudentID,
		StudentGroup: memberEvent.StudentGroupID,
	}

	switch memberEvent.Event {
	case events.EntityCreated:
		s.Add(ctx, trueM)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, memberEvent.ID, trueM)
	case events.EntityDeleted:
		s.Delete(ctx, memberEvent.ID)
	}
}

func (s *groupMemberService) Seed(ctx context.Context) error {
	members := []entities.GroupMember{}
	if err := jsonutil.ReadFileTo(filepath.Join("data", "group_members.json"), &members); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load group members seed data: %w", err), logger.ServiceDataFetchFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d group members updated successfully", len(members)),
		logger.ServiceOperationSuccess,
	)

	return nil
}

func (s *groupMemberService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	member entities.GroupMember,
) (entities.GroupMember, error) {
	updatedM, err := s.repository.ExternalUpdate(ctx, id, member)
	if err != nil {
		return entities.GroupMember{}, s.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"ExternalUpdate",
			err,
			logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedM),
		logger.ServiceOperationSuccess,
	)

	return updatedM, nil
}
