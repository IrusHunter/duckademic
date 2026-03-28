package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/student_group/entities"
	"github.com/IrusHunter/duckademic/services/student_group/repositories"
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

func NewGroupMemberService(
	gmr repositories.GroupMemberRepository,
	sr repositories.StudentRepository,
	gcr repositories.GroupCohortRepository,
	sgr repositories.StudentGroupRepository,
	eb events.EventBus,
) GroupMemberService {
	sc := platform.NewServiceConfig("GroupMembersService", filepath.Join("data", "group_members.json"), "group_member")

	res := &groupMemberService{
		repository:             gmr,
		studentRepository:      sr,
		groupCohortRepository:  gcr,
		studentGroupRepository: sgr,
		eventBus:               eb,
	}
	res.BaseService = platform.NewBaseServiceWithEventBus(sc, gmr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.GroupMember]{
			platform.OnAddPrepare: res.onAddPrepare,
		},
		eb,
	)
	res.logger = res.GetLogger()

	return res
}

type groupMemberService struct {
	platform.BaseService[entities.GroupMember]
	repository             repositories.GroupMemberRepository
	studentRepository      repositories.StudentRepository
	groupCohortRepository  repositories.GroupCohortRepository
	studentGroupRepository repositories.StudentGroupRepository
	logger                 logger.Logger
	eventBus               events.EventBus
}

func (s *groupMemberService) onAddPrepare(ctc context.Context, groupMember *entities.GroupMember) error {
	groupMember.ID = uuid.New()
	return nil
}

func (s *groupMemberService) Seed(ctx context.Context) error {
	groupMembersData := []struct {
		StudentSlug      string  `json:"student_slug"`
		GroupCohortName  string  `json:"group_cohort_name"`
		StudentGroupName *string `json:"student_group_name,omitempty"`
	}{}

	if err := jsonutil.ReadFileTo(filepath.Join("data", "group_members.json"), &groupMembersData); err != nil {
		return s.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"Seed",
			fmt.Errorf("failed to load group members seed data: %w", err),
			logger.ServiceValidationFailed,
		)
	}

	var lastError error
	for _, item := range groupMembersData {
		student := s.studentRepository.FindBySlug(ctx, item.StudentSlug)
		if student == nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("student with slug %q not found", item.StudentSlug),
				logger.ServiceValidationFailed,
			)
			continue
		}

		cohort := s.groupCohortRepository.FindFirstByName(ctx, item.GroupCohortName)
		if cohort == nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("group cohort with name %q not found", item.GroupCohortName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		var studentGroupID *uuid.UUID
		if item.StudentGroupName != nil {
			group := s.studentGroupRepository.FindFirstByName(ctx, *item.StudentGroupName)
			if group == nil {
				lastError = s.logger.LogAndReturnError(
					contextutil.GetTraceID(ctx),
					"Seed",
					fmt.Errorf("student group with name %q not found", *item.StudentGroupName),
					logger.ServiceValidationFailed,
				)
				continue
			}
			studentGroupID = &group.ID
		}

		gm := entities.GroupMember{
			ID:           uuid.New(),
			StudentID:    student.ID,
			GroupCohort:  cohort.ID,
			StudentGroup: studentGroupID,
		}

		_, err := s.Add(ctx, gm)
		if err != nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("failed to add group member %q: %w", gm.StudentID, err),
				logger.ServiceValidationFailed,
			)
		}
	}

	s.logger.Log(
		contextutil.GetTraceID(ctx),
		"Seed",
		fmt.Sprintf("%d group members processed", len(groupMembersData)),
		logger.ServiceOperationSuccess,
	)

	return lastError
}

func (s *groupMemberService) Add(
	ctx context.Context, member entities.GroupMember,
) (entities.GroupMember, error) {
	added, err := s.BaseService.Add(ctx, member)
	if err == nil {
		s.sendChanges(ctx, added, events.EntityCreated)
	}
	return added, err
}
func (s *groupMemberService) Delete(
	ctx context.Context, id uuid.UUID,
) (entities.GroupMember, error) {
	deleted, err := s.BaseService.Delete(ctx, id)
	if err == nil {
		s.sendChanges(ctx, deleted, events.EntityDeleted)
	}
	return deleted, err
}
func (s *groupMemberService) Update(
	ctx context.Context, id uuid.UUID, member entities.GroupMember,
) (entities.GroupMember, error) {
	updated, err := s.BaseService.Update(ctx, id, member)
	if err == nil {
		s.sendChanges(ctx, updated, events.EntityUpdated)
	}
	return updated, err
}

func (s *groupMemberService) sendChanges(
	ctx context.Context,
	member entities.GroupMember,
	eventType events.EventType,
) {
	eventM := events.GroupMemberRE{
		Event:          eventType,
		ID:             member.ID,
		StudentID:      member.StudentID,
		GroupCohortID:  member.GroupCohort,
		StudentGroupID: member.StudentGroup,
	}

	s.BaseService.SendChanges(ctx, eventM, eventType, events.GroupMemberRT)
}
