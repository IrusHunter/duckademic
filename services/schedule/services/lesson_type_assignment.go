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

type LessonTypeAssignmentService interface {
	platform.BaseService[entities.LessonTypeAssignment]
}

func NewLessonTypeAssignmentService(
	ltar repositories.LessonTypeAssignmentRepository,
	ltr repositories.LessonTypeRepository,
	dr repositories.DisciplineRepository,
	eb events.EventBus,
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
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.LessonTypeAssignment]{},
	)
	res.logger = res.GetLogger()

	eb.Subscribe(contextutil.SetTraceID(context.Background()), string(events.LessonTypeAssignmentRT), res.eventHandler)

	return res
}

type lessonTypeAssignmentService struct {
	platform.BaseService[entities.LessonTypeAssignment]
	repository           repositories.LessonTypeAssignmentRepository
	lessonTypeRepository repositories.LessonTypeRepository
	disciplineRepository repositories.DisciplineRepository
	logger               logger.Logger
}

func (s *lessonTypeAssignmentService) eventHandler(ctx context.Context, b []byte) {
	ltaEvent, err := events.FromByteConvertor[events.LessonTypeAssignmentRE](b)
	if err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "LessonTypeAssignmentRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "LessonTypeAssignmentRTHandlerRTHandler",
		fmt.Sprintf("received %s", ltaEvent), logger.EventDataReceived,
	)

	trueLTA := entities.LessonTypeAssignment{
		ID:            ltaEvent.ID,
		LessonTypeID:  ltaEvent.LessonTypeID,
		DisciplineID:  ltaEvent.DisciplineID,
		RequiredHours: ltaEvent.RequiredHours,
	}

	switch ltaEvent.Event {
	case events.EntityCreated:
		s.Add(ctx, trueLTA)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, ltaEvent.ID, trueLTA)
	case events.EntityDeleted:
		s.Delete(ctx, ltaEvent.ID)
	}
}

func (s *lessonTypeAssignmentService) Seed(ctx context.Context) error {
	assignments := []struct {
		LessonTypeName string `json:"lesson_type_name"`
		DisciplineName string `json:"discipline_name"`
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

		lta := s.repository.FindByLessonTypeAndDiscipline(ctx, lessonType.ID, discipline.ID)
		if lta == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("lta for %s and %s not found", discipline, lessonType), logger.ServiceValidationFailed,
			)
			continue
		}

		trueLta := entities.LessonTypeAssignment{}

		_, err := s.Update(ctx, lta.ID, trueLta)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to update %s: %w", lta, err), logger.ServiceValidationFailed,
			)
			continue
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d lesson type assignments processed from seed", len(assignments)), logger.ServiceOperationSuccess,
	)

	return lastError
}

func (s *lessonTypeAssignmentService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	lta entities.LessonTypeAssignment,
) (entities.LessonTypeAssignment, error) {
	updated, err := s.repository.ExternalUpdate(ctx, id, lta)
	if err != nil {
		return entities.LessonTypeAssignment{}, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "ExternalUpdate",
			err, logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%v successfully updated", updated), logger.ServiceOperationSuccess)
	return updated, nil
}
