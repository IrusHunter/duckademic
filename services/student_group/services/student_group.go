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

type StudentGroupService interface {
	platform.BaseService[entities.StudentGroup]
}

func NewStudentGroupService(
	sgr repositories.StudentGroupRepository,
	gcr repositories.GroupCohortRepository,
) StudentGroupService {
	sc := platform.NewServiceConfig("StudentGroupService", filepath.Join("data", "student_groups.json"), "student_group")

	res := &studentGroupService{
		repository: sgr,
	}
	res.BaseService = platform.NewBaseService(sc, sgr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.StudentGroup]{
			platform.OnAddPrepare:   res.onAddPrepare,
			platform.ValidateEntity: res.validateEntity,
		},
	)
	res.logger = res.GetLogger()

	return res
}

type studentGroupService struct {
	platform.BaseService[entities.StudentGroup]
	repository      repositories.StudentGroupRepository
	groupCohortRepo repositories.GroupCohortRepository
	logger          logger.Logger
}

func (s *studentGroupService) validateEntity(ctx context.Context, sg *entities.StudentGroup) error {
	if err := sg.ValidateName(); err != nil {
		return err
	}

	// Optionally: check that GroupCohortID is set
	if sg.GroupCohortID == uuid.Nil {
		return fmt.Errorf("group cohort ID required")
	}

	return nil
}

func (s *studentGroupService) onAddPrepare(ctx context.Context, sg *entities.StudentGroup) error {
	slugStr := slug.Make(sg.Name)
	if other := s.repository.FindBySlug(ctx, slugStr); other != nil {
		return fmt.Errorf("student group with slug %q already exists", slugStr)
	}
	sg.ID = uuid.New()
	sg.Slug = slugStr

	return nil
}

func (s *studentGroupService) Seed(ctx context.Context) error {
	studentGroupsData := []struct {
		Name            string `json:"name"`
		GroupCohortName string `json:"group_cohort_name"`
	}{}

	if err := jsonutil.ReadFileTo(filepath.Join("data", "student_groups.json"), &studentGroupsData); err != nil {
		return s.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"Seed",
			fmt.Errorf("failed to load student groups seed data: %w", err),
			logger.ServiceValidationFailed,
		)
	}

	var lastError error
	for _, item := range studentGroupsData {
		cohort := s.groupCohortRepo.FindFirstByName(ctx, item.GroupCohortName)
		if cohort == nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("group cohort with name %q not found", item.GroupCohortName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		sg := entities.StudentGroup{
			Name:          item.Name,
			GroupCohortID: cohort.ID,
		}

		_, err := s.Add(ctx, sg)
		if err != nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("failed to add student group %q: %w", sg.Name, err),
				logger.ServiceValidationFailed,
			)
		}
	}

	s.logger.Log(
		contextutil.GetTraceID(ctx),
		"Seed",
		fmt.Sprintf("%d student groups processed", len(studentGroupsData)),
		logger.ServiceOperationSuccess,
	)

	return lastError
}
