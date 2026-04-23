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

type LessonTypeAssignmentService interface {
	platform.BaseService[entities.LessonTypeAssignment]
	GetByDisciplineIDs(context.Context, []uuid.UUID) ([]entities.LessonTypeAssignment, error)
	GetUniqueLessonTypeIDs([]entities.LessonTypeAssignment) []uuid.UUID
	ToGeneratorLessonTypeAssignments(context.Context, []entities.LessonTypeAssignment) []GeneratorLessonTypeAssignment
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

func (s *lessonTypeAssignmentService) GetByDisciplineIDs(
	ctx context.Context,
	disciplineIDs []uuid.UUID,
) ([]entities.LessonTypeAssignment, error) {
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup
	var mu sync.Mutex

	var result []entities.LessonTypeAssignment
	var lastError error

	for _, disciplineID := range disciplineIDs {
		wg.Add(1)
		sem <- struct{}{}

		go func(disciplineID uuid.UUID) {
			defer wg.Done()
			defer func() { <-sem }()

			assignments, err := s.repository.GetByDisciplineID(ctx, disciplineID)
			if err != nil {
				mu.Lock()
				lastError = s.GetLogger().LogAndReturnError(
					contextutil.GetTraceID(ctx),
					"GetByDisciplineIDs",
					err,
					logger.ServiceRepositoryFailed,
				)
				mu.Unlock()
				return
			}

			mu.Lock()
			result = append(result, assignments...)
			mu.Unlock()

		}(disciplineID)
	}

	wg.Wait()

	return result, lastError
}
func (s *lessonTypeAssignmentService) GetUniqueLessonTypeIDs(lta []entities.LessonTypeAssignment) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{})
	result := make([]uuid.UUID, 0)

	for _, item := range lta {
		if _, ok := seen[item.LessonTypeID]; !ok {
			seen[item.LessonTypeID] = struct{}{}
			result = append(result, item.LessonTypeID)
		}
	}

	return result
}

type GeneratorLessonTypeAssignment struct {
	ID            uuid.UUID `json:"id"`
	LessonTypeID  uuid.UUID `json:"lesson_type_id"`
	DisciplineID  uuid.UUID `json:"discipline_id"`
	RequiredHours int       `json:"required_hours"`
}

func (s *lessonTypeAssignmentService) ToGeneratorLessonTypeAssignments(
	ctx context.Context,
	lta []entities.LessonTypeAssignment,
) []GeneratorLessonTypeAssignment {
	res := make([]GeneratorLessonTypeAssignment, 0, len(lta))

	for _, lessonTypeAssignment := range lta {
		res = append(res, GeneratorLessonTypeAssignment{
			ID:            lessonTypeAssignment.ID,
			LessonTypeID:  lessonTypeAssignment.LessonTypeID,
			DisciplineID:  lessonTypeAssignment.DisciplineID,
			RequiredHours: lessonTypeAssignment.RequiredHours,
		})
	}

	return res
}
