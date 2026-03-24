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

type LessonTypeAssignmentService interface {
	platform.BaseService[entities.LessonTypeAssignment]
}

func NewLessonTypeAssignmentService(
	ltar repositories.LessonTypeAssignmentRepository,
	ltr repositories.LessonTypeRepository,
	dr repositories.DisciplineRepository,
) LessonTypeAssignmentService {
	sc := platform.NewServiceConfig(
		"LessonTypeAssignmentService",
		filepath.Join("data", "lesson_type_assignments.json"),
		"lesson_type_assignment",
	)

	res := &lessonTypeAssignmentService{
		repository:           ltar,
		lessonTypeRepository: ltr,
		disciplineRepository: dr,
	}

	res.BaseService = platform.NewBaseService(sc, ltar,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.LessonTypeAssignment]{
			platform.OnAddPrepare:   res.onAddPrepare,
			platform.ValidateEntity: res.validateEntity,
		},
	)

	res.logger = res.GetLogger()
	return res
}

type lessonTypeAssignmentService struct {
	platform.BaseService[entities.LessonTypeAssignment]
	repository           repositories.LessonTypeAssignmentRepository
	lessonTypeRepository repositories.LessonTypeRepository
	disciplineRepository repositories.DisciplineRepository
	logger               logger.Logger
}

func (s *lessonTypeAssignmentService) validateEntity(ctx context.Context, lta *entities.LessonTypeAssignment) error {
	if err := lta.ValidateRequiredHours(); err != nil {
		return err
	}
	return nil
}
func (s *lessonTypeAssignmentService) onAddPrepare(ctx context.Context, lta *entities.LessonTypeAssignment) error {
	if existing := s.repository.FindByLessonTypeAndDiscipline(ctx, lta.LessonTypeID, lta.DisciplineID); existing != nil {
		return fmt.Errorf("assignment for lesson type with id %q and discipline with id %q already exists",
			lta.LessonTypeID, lta.DisciplineID,
		)
	}

	lta.ID = uuid.New()
	return nil
}

func (s *lessonTypeAssignmentService) Seed(ctx context.Context) error {
	assignments := []struct {
		LessonTypeName string `json:"lesson_type_name"`
		DisciplineName string `json:"discipline_name"`
		RequiredHours  int    `json:"required_hours"`
	}{}

	if err := jsonutil.ReadFileTo(filepath.Join("data", "lesson_type_assignments.json"), &assignments); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load lesson type assignments seed data: %w", err), logger.ServiceValidationFailed,
		)
	}

	var lastError error
	for _, item := range assignments {
		lessonType := s.lessonTypeRepository.FindFirstByName(ctx, item.LessonTypeName)
		if lessonType == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("lesson type %q not found", item.LessonTypeName), logger.ServiceValidationFailed,
			)
			continue
		}

		discipline := s.disciplineRepository.FindFirstByName(ctx, item.DisciplineName)
		if discipline == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("discipline %q not found", item.DisciplineName), logger.ServiceValidationFailed,
			)
			continue
		}

		lta := entities.LessonTypeAssignment{
			LessonTypeID:  lessonType.ID,
			DisciplineID:  discipline.ID,
			RequiredHours: item.RequiredHours,
		}

		_, err := s.Add(ctx, lta)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to add %s: %w", lta, err), logger.ServiceValidationFailed,
			)
			continue
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d lesson type assignments processed from seed", len(assignments)), logger.ServiceOperationSuccess,
	)

	return lastError
}
