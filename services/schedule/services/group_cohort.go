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

type GroupCohortService interface {
	platform.BaseService[entities.GroupCohort]
	GetBySemesterIDs(context.Context, []uuid.UUID) ([]entities.GroupCohort, error)
	ToGeneratorGroupCohorts(context.Context, []entities.GroupCohort) []GeneratorGroupCohort
	GetUniqueGroupCohortIDs([]entities.GroupCohort) []uuid.UUID
}

func NewGroupCohortService(
	gr repositories.GroupCohortRepository,
	sgr repositories.StudentGroupRepository,
	gms GroupMemberService,
	eb events.EventBus,
) GroupCohortService {
	sc := platform.NewServiceConfig(
		"GroupCohortService",
		filepath.Join("data", "group_cohorts.json"),
		"group cohort",
	)

	res := &groupCohortService{
		repository:             gr,
		studentGroupREpository: sgr,
		groupMembersService:    gms,
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
	repository             repositories.GroupCohortRepository
	studentGroupREpository repositories.StudentGroupRepository
	groupMembersService    GroupMemberService
	logger                 logger.Logger
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
		ID:         cohortEvent.ID,
		Slug:       cohortEvent.Slug,
		Name:       cohortEvent.Name,
		SemesterID: &cohortEvent.SemesterID,
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
func (s *groupCohortService) GetBySemesterIDs(ctx context.Context, semesterIDs []uuid.UUID) ([]entities.GroupCohort, error) {
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup
	var mu sync.Mutex

	var result []entities.GroupCohort
	var lastError error

	for _, semesterID := range semesterIDs {
		wg.Add(1)
		sem <- struct{}{}

		go func(semesterID uuid.UUID) {
			defer wg.Done()
			defer func() { <-sem }()

			cohorts, err := s.repository.GetBySemesterID(ctx, semesterID)
			if err != nil {
				mu.Lock()
				lastError = s.GetLogger().LogAndReturnError(
					contextutil.GetTraceID(ctx),
					"GetBySemesterIDs",
					err,
					logger.ServiceRepositoryFailed,
				)
				mu.Unlock()
				return
			}

			mu.Lock()
			result = append(result, cohorts...)
			mu.Unlock()

		}(semesterID)
	}

	wg.Wait()

	return result, lastError
}

type GeneratorGroupCohort struct {
	ID     uuid.UUID               `json:"id"`
	Name   string                  `json:"name"`
	Groups []GeneratorStudentGroup `json:"groups"`
}

func (s *groupCohortService) ToGeneratorGroupCohorts(ctx context.Context, gc []entities.GroupCohort) []GeneratorGroupCohort {
	res := make([]GeneratorGroupCohort, 0, len(gc))

	for _, groupCohort := range gc {
		studentGroups, err := s.studentGroupREpository.GetByGroupCohortID(ctx, groupCohort.ID)
		if err != nil {
			s.GetLogger().LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"ToGeneratorGroupCohorts",
				err,
				logger.ServiceRepositoryFailed,
			)
			continue
		}

		generatorSGs := s.groupMembersService.ToGeneratorStudentGroup(ctx, studentGroups)
		res = append(res, GeneratorGroupCohort{
			ID:     groupCohort.ID,
			Name:   groupCohort.Name,
			Groups: generatorSGs,
		})
	}

	return res
}
func (s *groupCohortService) GetUniqueGroupCohortIDs(cohorts []entities.GroupCohort) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{})
	result := make([]uuid.UUID, 0)

	for _, c := range cohorts {
		if _, ok := seen[c.ID]; !ok {
			seen[c.ID] = struct{}{}
			result = append(result, c.ID)
		}
	}

	return result
}
