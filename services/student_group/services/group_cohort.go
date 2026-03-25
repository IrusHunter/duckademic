package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/student_group/entities"
	"github.com/IrusHunter/duckademic/services/student_group/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type GroupCohortService interface {
	platform.BaseService[entities.GroupCohort]
}

func NewGroupCohortService(
	gcr repositories.GroupCohortRepository,
	sr repositories.SemesterRepository,
) GroupCohortService {
	sc := platform.NewServiceConfig("GroupCohortService", filepath.Join("data", "group_cohorts.json"), "group_cohort")

	res := &groupCohortService{
		repository: gcr,
	}
	res.BaseService = platform.NewBaseService(sc, gcr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.GroupCohort]{
			platform.OnAddPrepare:   res.onAddPrepare,
			platform.ValidateEntity: res.validateEntity,
		},
	)
	res.logger = res.GetLogger()

	return res
}

type groupCohortService struct {
	platform.BaseService[entities.GroupCohort]
	repository         repositories.GroupCohortRepository
	semesterRepository repositories.SemesterRepository
	logger             logger.Logger
}

func (s *groupCohortService) validateEntity(ctx context.Context, groupCohort *entities.GroupCohort) error {
	if err := groupCohort.ValidateName(); err != nil {
		return err
	}

	return nil
}
func (s *groupCohortService) onAddPrepare(ctx context.Context, groupCohort *entities.GroupCohort) error {
	slug := slug.Make(groupCohort.Name)
	if other := s.repository.FindBySlug(ctx, slug); other != nil {
		return fmt.Errorf("group cohort with slug %q already exists", slug)
	}
	groupCohort.ID = uuid.New()
	groupCohort.Slug = slug

	return nil
}

func (s *groupCohortService) Seed(ctx context.Context) error {
	groupCohortsData := []struct {
		Name           string `json:"name"`
		CurriculumName string `json:"curriculum_name"`
		SemesterNumber int    `json:"semester_number"`
	}{}

	if err := jsonutil.ReadFileTo(filepath.Join("data", "group_cohorts.json"), &groupCohortsData); err != nil {
		return s.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"Seed",
			fmt.Errorf("failed to load group cohorts seed data: %w", err),
			logger.ServiceValidationFailed,
		)
	}

	var lastError error
	for _, item := range groupCohortsData {
		semesterSlug := fmt.Sprintf("%s-%d", slug.Make(item.CurriculumName), item.SemesterNumber)
		semester := s.semesterRepository.FindBySlug(ctx, semesterSlug)
		if semester == nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("semester slug %q not found", semesterSlug),
				logger.ServiceValidationFailed,
			)
			continue
		}

		gc := entities.GroupCohort{
			Name:       item.Name,
			SemesterID: semester.ID,
		}

		_, err := s.Add(ctx, gc)
		if err != nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("failed to add group cohort %q: %w", gc.Name, err),
				logger.ServiceValidationFailed,
			)
		}
	}

	s.logger.Log(
		contextutil.GetTraceID(ctx),
		"Seed",
		fmt.Sprintf("%d group cohorts processed", len(groupCohortsData)),
		logger.ServiceOperationSuccess,
	)

	return lastError
}
