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

type LessonTypeService interface {
	platform.BaseService[entities.LessonType]
}

func NewLessonTypeService(ltr repositories.LessonTypeRepository, eb events.EventBus) LessonTypeService {
	sc := platform.NewServiceConfig("LessonTypeService", filepath.Join("data", "lesson_types.json"), "lesson_type")

	res := &lessonTypeService{
		repository: ltr,
	}
	res.BaseService = platform.NewBaseService(sc, ltr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.LessonType]{},
	)
	res.logger = res.GetLogger()

	eb.Subscribe(contextutil.SetTraceID(context.Background()), string(events.LessonTypeRT), res.eventHandler)

	return res
}

type lessonTypeService struct {
	platform.BaseService[entities.LessonType]
	repository repositories.LessonTypeRepository
	logger     logger.Logger
}

func (s *lessonTypeService) eventHandler(ctx context.Context, b []byte) {
	ltr, err := events.FromByteConvertor[events.LessonTypeRE](b)
	if err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "LessonTypeRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "LessonTypeRTHandler",
		fmt.Sprintf("received %s", ltr), logger.EventDataReceived,
	)

	trueLTR := entities.LessonType{
		ID:   ltr.ID,
		Name: ltr.Name,
		Slug: ltr.Slug,
	}

	switch ltr.Event {
	case events.EntityCreated:
		s.Add(ctx, trueLTR)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, ltr.ID, trueLTR)
	case events.EntityDeleted:
		s.Delete(ctx, ltr.ID)
	}
}

func (s *lessonTypeService) Seed(ctx context.Context) error {
	lessonTypes := []entities.LessonType{}
	if err := jsonutil.ReadFileTo(filepath.Join("data", "lesson_types.json"), &lessonTypes); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load lesson_types seed data: %w", err), logger.ServiceDataFetchFailed,
		)
	}

	var lastError error
	for _, lt := range lessonTypes {
		tlt := s.repository.FindFirstByName(ctx, lt.Name)
		if tlt == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("lesson type with name %q not found", lt.Name), logger.ServiceDataFetchFailed,
			)
			continue
		}

		_, err := s.Update(ctx, tlt.ID, lt)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to update %s: %w", lt, err), logger.ServiceRepositoryFailed,
			)
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d lesson types updated successfully", len(lessonTypes)), logger.ServiceOperationSuccess,
	)
	return lastError
}

func (s *lessonTypeService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	lessonType entities.LessonType,
) (entities.LessonType, error) {
	updatedLT, err := s.repository.ExternalUpdate(ctx, id, lessonType)
	if err != nil {
		return entities.LessonType{}, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "ExternalUpdate",
			err, logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedLT), logger.ServiceOperationSuccess)
	return updatedLT, nil
}
