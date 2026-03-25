package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/curriculum/entities"
	"github.com/IrusHunter/duckademic/services/curriculum/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type LessonTypeService interface {
	platform.BaseService[entities.LessonType]
}

func NewLessonTypeService(
	ltr repositories.LessonTypeRepository,
	eb events.EventBus,
) LessonTypeService {
	sc := platform.NewServiceConfig(
		"LessonTypeService",
		filepath.Join("data", "lesson_types.json"),
		"lesson_type",
	)

	res := &lessonTypeService{
		repository: ltr,
		eventBus:   eb,
	}

	res.BaseService = platform.NewBaseServiceWithEventBus(sc, ltr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.LessonType]{
			platform.OnAddPrepare:   res.onAddPrepare,
			platform.ValidateEntity: res.validateEntity,
		},
		eb,
	)

	res.logger = res.GetLogger()

	return res
}

type lessonTypeService struct {
	platform.BaseService[entities.LessonType]
	repository repositories.LessonTypeRepository
	logger     logger.Logger
	eventBus   events.EventBus
}

func (s *lessonTypeService) validateEntity(ctx context.Context, lt *entities.LessonType) error {
	if err := lt.ValidateName(); err != nil {
		return err
	}
	if err := lt.ValidateHoursValue(); err != nil {
		return err
	}

	return nil
}
func (s *lessonTypeService) onAddPrepare(ctx context.Context, lt *entities.LessonType) error {
	slugStr := slug.Make(lt.Name)
	if other := s.repository.FindBySlug(ctx, slugStr); other != nil {
		return fmt.Errorf("lesson type with slug %q already exists", slugStr)
	}
	lt.ID = uuid.New()
	lt.Slug = slugStr
	return nil
}

func (s *lessonTypeService) Seed(ctx context.Context) error {
	lessonTypes := []entities.LessonType{}
	if err := jsonutil.ReadFileTo(filepath.Join("data", "lesson_types.json"), &lessonTypes); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load lesson types seed data: %w", err), logger.ServiceValidationFailed,
		)
	}

	var lastError error
	for _, lessonType := range lessonTypes {
		_, err := s.Add(ctx, lessonType)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to add %s: %w", lessonType.String(), err), logger.ServiceValidationFailed,
			)
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d lesson types added successfully", len(lessonTypes)), logger.ServiceOperationSuccess,
	)

	return lastError
}

func (s *lessonTypeService) Add(
	ctx context.Context, lessonType entities.LessonType,
) (entities.LessonType, error) {
	addedLT, err := s.BaseService.Add(ctx, lessonType)
	if err == nil {
		s.sendChanges(ctx, addedLT, events.EntityCreated)
	}
	return addedLT, err
}
func (s *lessonTypeService) Delete(
	ctx context.Context, id uuid.UUID,
) (entities.LessonType, error) {
	deletedLT, err := s.BaseService.Delete(ctx, id)
	if err == nil {
		s.sendChanges(ctx, deletedLT, events.EntityDeleted)
	}
	return deletedLT, err
}
func (s *lessonTypeService) Update(
	ctx context.Context, id uuid.UUID, lessonType entities.LessonType,
) (entities.LessonType, error) {
	updatedLT, err := s.BaseService.Update(ctx, id, lessonType)
	if err == nil {
		s.sendChanges(ctx, updatedLT, events.EntityUpdated)
	}
	return updatedLT, err
}

func (s *lessonTypeService) sendChanges(
	ctx context.Context,
	lessonType entities.LessonType,
	event events.EventType,
) {
	eventLT := events.LessonTypeRE{
		Event:      event,
		ID:         lessonType.ID,
		Slug:       lessonType.Slug,
		Name:       lessonType.Name,
		HoursValue: lessonType.HoursValue,
	}

	s.BaseService.SendChanges(ctx, eventLT, event, events.LessonTypeRT)
}
