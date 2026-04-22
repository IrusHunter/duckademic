package services

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type GroupMemberService interface {
	platform.BaseService[entities.GroupMember]
	ToGeneratorStudentGroup(context.Context, []entities.StudentGroup) []GeneratorStudentGroup
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
		ID:             memberEvent.ID,
		StudentID:      memberEvent.StudentID,
		StudentGroupID: memberEvent.StudentGroupID,
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

type GeneratorStudentGroup struct {
	ID              uuid.UUID   `json:"id"`
	Name            string      `json:"name"`
	StudentCount    int         `json:"student_count"`
	ConnectedGroups []uuid.UUID `json:"connected_groups"`
}

func (s *groupMemberService) ToGeneratorStudentGroup(ctx context.Context, sg []entities.StudentGroup) []GeneratorStudentGroup {
	res := make([]GeneratorStudentGroup, 0, len(sg))

	traceID := contextutil.GetTraceID(ctx)

	for _, studentGroup := range sg {
		studentsIDs, err := s.repository.GetByGroupID(ctx, studentGroup.ID)
		if err != nil {
			s.GetLogger().LogAndReturnError(
				traceID,
				"ToGeneratorStudentGroup:GetByGroupID",
				err,
				logger.ServiceRepositoryFailed,
			)
			studentsIDs = []uuid.UUID{}
		}

		allConnectedGMs, err := s.repository.GetByStudentIDs(ctx, studentsIDs)
		if err != nil {
			s.GetLogger().LogAndReturnError(
				traceID,
				"ToGeneratorStudentGroup:GetByStudentIDs",
				err,
				logger.ServiceRepositoryFailed,
			)
			allConnectedGMs = append(allConnectedGMs, entities.GroupMember{StudentGroupID: &studentGroup.ID})
		}

		uniqueSGids := s.GetUniqueStudentGroupIDs(allConnectedGMs)
		ind := slices.Index(uniqueSGids, studentGroup.ID)
		if ind != -1 {
			uniqueSGids = append(uniqueSGids[:ind], uniqueSGids[ind+1:]...)
		}

		res = append(res, GeneratorStudentGroup{
			ID:              studentGroup.ID,
			Name:            studentGroup.Name,
			StudentCount:    len(studentsIDs),
			ConnectedGroups: uniqueSGids,
		})
	}

	return res
}
func (s *groupMemberService) GetUniqueStudentGroupIDs(gm []entities.GroupMember) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{})
	result := make([]uuid.UUID, 0)

	for _, item := range gm {
		if item.StudentGroupID == nil {
			continue
		}
		if _, ok := seen[*item.StudentGroupID]; !ok {
			seen[*item.StudentGroupID] = struct{}{}
			result = append(result, *item.StudentGroupID)
		}
	}

	return result
}
