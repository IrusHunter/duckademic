package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/asset/entities"
	"github.com/IrusHunter/duckademic/services/asset/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type ClassroomService interface {
	platform.BaseService[entities.Classroom]
}

func NewClassroomService(
	cr repositories.ClassroomRepository,
	eb events.EventBus,
) ClassroomService {
	sc := platform.NewServiceConfig(
		"ClassroomService",
		filepath.Join("data", "classrooms.json"),
		"classroom",
	)

	res := &classroomService{
		repository: cr,
		eventBus:   eb,
	}

	res.BaseService = platform.NewBaseServiceWithEventBus(sc, cr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Classroom]{
			platform.OnAddPrepare:   res.onAddPrepare,
			platform.ValidateEntity: res.validateEntity,
		},
		eb,
	)

	res.logger = res.GetLogger()

	return res
}

type classroomService struct {
	platform.BaseService[entities.Classroom]
	repository repositories.ClassroomRepository
	logger     logger.Logger
	eventBus   events.EventBus
}

func (s *classroomService) validateEntity(ctx context.Context, c *entities.Classroom) error {
	if err := c.ValidateNumber(); err != nil {
		return err
	}
	if err := c.ValidateCapacity(); err != nil {
		return err
	}
	return nil
}
func (s *classroomService) onAddPrepare(ctx context.Context, c *entities.Classroom) error {
	slugStr := slug.Make(c.Number)
	if other := s.repository.FindBySlug(ctx, slugStr); other != nil {
		return fmt.Errorf("classroom with slug %q already exists", slugStr)
	}
	c.ID = uuid.New()
	c.Slug = slugStr
	return nil
}

func (s *classroomService) Seed(ctx context.Context) error {
	classrooms := []entities.Classroom{}
	if err := jsonutil.ReadFileTo(filepath.Join("data", "classrooms.json"), &classrooms); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load classrooms seed data: %w", err), logger.ServiceValidationFailed,
		)
	}

	var lastError error
	for _, classroom := range classrooms {
		_, err := s.Add(ctx, classroom)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to add %s: %w", classroom.String(), err), logger.ServiceValidationFailed,
			)
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d classrooms added successfully", len(classrooms)), logger.ServiceOperationSuccess,
	)

	return lastError
}

func (s *classroomService) Add(
	ctx context.Context, classroom entities.Classroom,
) (entities.Classroom, error) {
	added, err := s.BaseService.Add(ctx, classroom)
	if err == nil {
		s.sendChanges(ctx, added, events.EntityCreated)
	}
	return added, err
}

func (s *classroomService) Delete(
	ctx context.Context, id uuid.UUID,
) (entities.Classroom, error) {
	deleted, err := s.BaseService.Delete(ctx, id)
	if err == nil {
		s.sendChanges(ctx, deleted, events.EntityDeleted)
	}
	return deleted, err
}

func (s *classroomService) Update(
	ctx context.Context, id uuid.UUID, classroom entities.Classroom,
) (entities.Classroom, error) {
	updated, err := s.BaseService.Update(ctx, id, classroom)
	if err == nil {
		s.sendChanges(ctx, updated, events.EntityUpdated)
	}
	return updated, err
}

func (s *classroomService) sendChanges(
	ctx context.Context,
	classroom entities.Classroom,
	event events.EventType,
) {
	eventC := events.ClassroomRE{
		Event:    event,
		ID:       classroom.ID,
		Slug:     classroom.Slug,
		Number:   classroom.Number,
		Capacity: classroom.Capacity,
	}

	s.BaseService.SendChanges(ctx, eventC, event, events.ClassroomRT)
}
