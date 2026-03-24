package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/curriculum/entities"
	"github.com/IrusHunter/duckademic/services/curriculum/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type SemesterService interface {
	platform.BaseService[entities.Semester]
}

func NewSemesterService(
	sr repositories.SemesterRepository,
	cr repositories.CurriculumRepository,
) SemesterService {
	sc := platform.NewServiceConfig(
		"SemesterService",
		filepath.Join("data", "semesters.json"),
		"semester",
	)

	res := &semesterService{
		repository:           sr,
		curriculumRepository: cr,
	}

	res.BaseService = platform.NewBaseService(sc, sr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Semester]{
			platform.OnAddPrepare:   res.onAddPrepare,
			platform.ValidateEntity: res.validateEntity,
		},
	)

	res.logger = res.GetLogger()

	return res
}

type semesterService struct {
	platform.BaseService[entities.Semester]
	repository           repositories.SemesterRepository
	curriculumRepository repositories.CurriculumRepository
	logger               logger.Logger
}

func (s *semesterService) validateEntity(ctx context.Context, semester *entities.Semester) error {
	if err := semester.ValidateNumber(); err != nil {
		return err
	}
	return nil
}
func (s *semesterService) onAddPrepare(ctx context.Context, semester *entities.Semester) error {
	curriculum := s.curriculumRepository.FindByID(ctx, semester.CurriculumID)
	if curriculum == nil {
		return fmt.Errorf("curriculum with id %s not found", semester.CurriculumID)
	}

	slugStr := fmt.Sprintf("%s-%d", curriculum.Slug, semester.Number)
	if other := s.repository.FindBySlug(ctx, slugStr); other != nil {
		return fmt.Errorf("semester with slug %q already exists", slugStr)
	}
	semester.ID = uuid.New()
	semester.Slug = slugStr
	return nil
}

func (s *semesterService) Seed(ctx context.Context) error {
	semesters := []struct {
		CurriculumName string `json:"curriculum_name"`
		Number         int    `json:"number"`
	}{}

	if err := jsonutil.ReadFileTo(filepath.Join("data", "semesters.json"), &semesters); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load semesters seed data: %w", err), logger.ServiceValidationFailed,
		)
	}

	var lastError error
	for _, semester := range semesters {
		curriculum := s.curriculumRepository.FindFirstByName(ctx, semester.CurriculumName)
		if curriculum == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("curriculum %q not found", semester.CurriculumName), logger.ServiceValidationFailed,
			)
			continue
		}

		trueSemester := entities.Semester{
			CurriculumID: curriculum.ID,
			Number:       semester.Number,
		}

		_, err := s.Add(ctx, trueSemester)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to add %s: %w", trueSemester, err),
				logger.ServiceValidationFailed,
			)
			continue
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d semesters processed from seed", len(semesters)), logger.ServiceOperationSuccess,
	)

	return lastError
}
