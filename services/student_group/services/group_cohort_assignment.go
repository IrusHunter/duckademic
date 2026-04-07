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

type GroupCohortAssignmentService interface {
	platform.BaseService[entities.GroupCohortAssignment]
}

func NewGroupCohortAssignmentService(
	gcar repositories.GroupCohortAssignmentRepository,
	gcr repositories.GroupCohortRepository,
	dr repositories.DisciplineRepository,
	ltr repositories.LessonTypeRepository,
	eb events.EventBus,
) GroupCohortAssignmentService {
	sc := platform.NewServiceConfig(
		"GroupCohortAssignmentService",
		filepath.Join("data", "group_cohort_assignments.json"),
		"group_cohort_assignment",
	)

	res := &groupCohortAssignmentService{
		repository:            gcar,
		groupCohortRepository: gcr,
		disciplineRepository:  dr,
		lessonTypeRepository:  ltr,
		eventBus:              eb,
	}
	res.BaseService = platform.NewBaseServiceWithEventBus(sc, gcar,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.GroupCohortAssignment]{
			platform.OnAddPrepare: res.onAddPrepare,
		},
		eb,
	)
	res.logger = res.GetLogger()

	return res
}

type groupCohortAssignmentService struct {
	platform.BaseService[entities.GroupCohortAssignment]
	repository            repositories.GroupCohortAssignmentRepository
	groupCohortRepository repositories.GroupCohortRepository
	disciplineRepository  repositories.DisciplineRepository
	lessonTypeRepository  repositories.LessonTypeRepository
	logger                logger.Logger
	eventBus              events.EventBus
}

func (s *groupCohortAssignmentService) onAddPrepare(
	ctx context.Context, assignment *entities.GroupCohortAssignment,
) error {
	assignment.ID = uuid.New()
	return nil
}

func (s *groupCohortAssignmentService) Seed(ctx context.Context) error {
	seedData := []struct {
		GroupCohortName string `json:"group_cohort_name"`
		DisciplineName  string `json:"discipline_name"`
		LessonTypeName  string `json:"lesson_type_name"`
	}{}

	if err := jsonutil.ReadFileTo(
		filepath.Join("data", "group_cohort_assignments.json"),
		&seedData,
	); err != nil {
		return s.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"Seed",
			fmt.Errorf("failed to load group cohort assignments seed data: %w", err),
			logger.ServiceValidationFailed,
		)
	}

	var lastError error
	for _, item := range seedData {
		groupCohort := s.groupCohortRepository.FindFirstByName(ctx, item.GroupCohortName)
		if groupCohort == nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("group cohort %q not found", item.GroupCohortName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		discipline := s.disciplineRepository.FindFirstByName(ctx, item.DisciplineName)
		if discipline == nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("discipline %q not found", item.DisciplineName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		lessonType := s.lessonTypeRepository.FindFirstByName(ctx, item.LessonTypeName)
		if lessonType == nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("lesson type %q not found", item.LessonTypeName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		assignment := entities.GroupCohortAssignment{
			GroupCohortID: groupCohort.ID,
			DisciplineID:  discipline.ID,
			LessonTypeID:  lessonType.ID,
		}

		_, err := s.Add(ctx, assignment)
		if err != nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("failed to add assignment: %w", err),
				logger.ServiceValidationFailed,
			)
		}
	}

	s.logger.Log(
		contextutil.GetTraceID(ctx),
		"Seed",
		fmt.Sprintf("%d group cohort assignments processed", len(seedData)),
		logger.ServiceOperationSuccess,
	)

	return lastError
}

func (s *groupCohortAssignmentService) Add(
	ctx context.Context, assignment entities.GroupCohortAssignment,
) (entities.GroupCohortAssignment, error) {
	added, err := s.BaseService.Add(ctx, assignment)
	if err == nil {
		s.sendChanges(ctx, added, events.EntityCreated)
	}
	return added, err
}

func (s *groupCohortAssignmentService) Delete(
	ctx context.Context, id uuid.UUID,
) (entities.GroupCohortAssignment, error) {
	deleted, err := s.BaseService.Delete(ctx, id)
	if err == nil {
		s.sendChanges(ctx, deleted, events.EntityDeleted)
	}
	return deleted, err
}

func (s *groupCohortAssignmentService) Update(
	ctx context.Context, id uuid.UUID, assignment entities.GroupCohortAssignment,
) (entities.GroupCohortAssignment, error) {
	updated, err := s.BaseService.Update(ctx, id, assignment)
	if err == nil {
		s.sendChanges(ctx, updated, events.EntityUpdated)
	}
	return updated, err
}

func (s *groupCohortAssignmentService) sendChanges(
	ctx context.Context,
	assignment entities.GroupCohortAssignment,
	eventType events.EventType,
) {
	event := events.GroupCohortAssignmentRE{
		Event:         eventType,
		ID:            assignment.ID,
		GroupCohortID: assignment.GroupCohortID,
		DisciplineID:  assignment.DisciplineID,
		LessonTypeID:  assignment.LessonTypeID,
	}
	s.BaseService.SendChanges(ctx, event, eventType, events.GroupCohortAssignmentRT)
}
