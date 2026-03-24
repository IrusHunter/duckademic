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

type SemesterDisciplineService interface {
	platform.BaseService[entities.SemesterDiscipline]
}

func NewSemesterDisciplineService(
	sdr repositories.SemesterDisciplineRepository,
	sr repositories.SemesterRepository,
	dr repositories.DisciplineRepository,
	cr repositories.CurriculumRepository,
) SemesterDisciplineService {
	sc := platform.NewServiceConfig(
		"SemesterDisciplineService",
		filepath.Join("data", "semester_discipline.json"),
		"semester_discipline",
	)

	res := &semesterDisciplineService{
		repository:           sdr,
		semesterRepository:   sr,
		disciplineRepository: dr,
		curriculumRepository: cr,
	}

	res.BaseService = platform.NewBaseService(sc, sdr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.SemesterDiscipline]{
			platform.OnAddPrepare: res.onAddPrepare,
		},
	)

	res.logger = res.GetLogger()
	return res
}

type semesterDisciplineService struct {
	platform.BaseService[entities.SemesterDiscipline]
	repository           repositories.SemesterDisciplineRepository
	semesterRepository   repositories.SemesterRepository
	disciplineRepository repositories.DisciplineRepository
	curriculumRepository repositories.CurriculumRepository
	logger               logger.Logger
}

func (s *semesterDisciplineService) onAddPrepare(
	ctx context.Context, semesterDiscipline *entities.SemesterDiscipline,
) error {
	semesterDiscipline.ID = uuid.New()
	return nil
}

func (s *semesterDisciplineService) Seed(ctx context.Context) error {
	type seedItem struct {
		CurriculumName string `json:"curriculum_name"`
		SemesterNumber int    `json:"semester_number"`
		DisciplineName string `json:"discipline_name"`
	}

	var mappings []seedItem
	if err := jsonutil.ReadFileTo(filepath.Join("data", "semester_discipline.json"), &mappings); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load semester discipline seed data: %w", err), logger.ServiceValidationFailed,
		)
	}

	var lastError error
	curriculumCache := make(map[string]entities.Curriculum)
	semesterCache := make(map[string]entities.Semester)

	for _, item := range mappings {
		if _, ok := curriculumCache[item.CurriculumName]; !ok {
			curr := s.curriculumRepository.FindFirstByName(ctx, item.CurriculumName)
			if curr == nil {
				lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
					fmt.Errorf("curriculum %q not found", item.CurriculumName), logger.ServiceValidationFailed,
				)
				continue
			}
			curriculumCache[item.CurriculumName] = *curr
		}
		curriculum := curriculumCache[item.CurriculumName]

		semesterKey := fmt.Sprintf("%s-%d", curriculum.ID, item.SemesterNumber)
		if _, ok := semesterCache[semesterKey]; !ok {
			sem := s.semesterRepository.FindByCurriculumIDAndNumber(ctx, curriculum.ID, item.SemesterNumber)
			if sem == nil {
				lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
					fmt.Errorf("semester number %d not found for curriculum %q", item.SemesterNumber, item.CurriculumName),
					logger.ServiceValidationFailed,
				)
				continue
			}
			semesterCache[semesterKey] = *sem
		}
		semester := semesterCache[semesterKey]

		discipline := s.disciplineRepository.FindFirstByName(ctx, item.DisciplineName)
		if discipline == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("discipline %q not found", item.DisciplineName), logger.ServiceValidationFailed,
			)
			continue
		}

		sd := entities.SemesterDiscipline{
			SemesterID:   semester.ID,
			DisciplineID: discipline.ID,
		}

		_, err := s.Add(ctx, sd)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to add %s: %w", sd, err), logger.ServiceValidationFailed,
			)
			continue
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d semester discipline mappings processed from seed", len(mappings)), logger.ServiceOperationSuccess,
	)

	return lastError
}
