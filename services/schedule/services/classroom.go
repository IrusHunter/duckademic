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

type ClassroomService interface {
	platform.BaseService[entities.Classroom]
}

func NewClassroomService(cr repositories.ClassroomRepository, eb events.EventBus) ClassroomService {
	sc := platform.NewServiceConfig("ClassroomService", filepath.Join("data", "classrooms.json"), "classroom")

	res := &classroomService{
		repository: cr,
	}

	res.BaseService = platform.NewBaseService(
		sc,
		cr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Classroom]{},
	)

	res.logger = res.GetLogger()

	eb.Subscribe(
		contextutil.SetTraceID(context.Background()),
		string(events.ClassroomRT),
		res.eventHandler,
	)

	return res
}

type classroomService struct {
	platform.BaseService[entities.Classroom]
	repository repositories.ClassroomRepository
	logger     logger.Logger
}

func (s *classroomService) eventHandler(ctx context.Context, b []byte) {
	cr, err := events.FromByteConvertor[events.ClassroomRE](b)
	if err != nil {
		s.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"ClassroomRTHandler",
			err,
			logger.EventDataReadFailed,
		)
		return
	}

	s.logger.Log(
		contextutil.GetTraceID(ctx),
		"ClassroomRTHandler",
		fmt.Sprintf("received %s", cr),
		logger.EventDataReceived,
	)

	trueCR := entities.Classroom{
		ID:       cr.ID,
		Slug:     cr.Slug,
		Number:   cr.Number,
		Capacity: cr.Capacity,
	}

	switch cr.Event {
	case events.EntityCreated:
		s.Add(ctx, trueCR)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, cr.ID, trueCR)
	case events.EntityDeleted:
		s.Delete(ctx, cr.ID)
	}
}

func (s *classroomService) Seed(ctx context.Context) error {
	classrooms := []entities.Classroom{}
	if err := jsonutil.ReadFileTo(filepath.Join("data", "classrooms.json"), &classrooms); err != nil {
		return s.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"Seed",
			fmt.Errorf("failed to load classrooms seed data: %w", err),
			logger.ServiceDataFetchFailed,
		)
	}

	var lastError error
	for _, classroom := range classrooms {
		existing := s.repository.FindFirstByNumber(ctx, classroom.Number)
		if existing == nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("classroom with number %q not found", classroom.Number),
				logger.ServiceDataFetchFailed,
			)
			continue
		}

		_, err := s.Update(ctx, existing.ID, classroom)
		if err != nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("failed to update %v: %w", classroom, err),
				logger.ServiceRepositoryFailed,
			)
		}
	}

	s.logger.Log(
		contextutil.GetTraceID(ctx),
		"Seed",
		fmt.Sprintf("%d classrooms updated successfully", len(classrooms)),
		logger.ServiceOperationSuccess,
	)

	return lastError
}

func (s *classroomService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	classroom entities.Classroom,
) (entities.Classroom, error) {

	updatedCR, err := s.repository.ExternalUpdate(ctx, id, classroom)
	if err != nil {
		return entities.Classroom{}, s.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"ExternalUpdate",
			err,
			logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(
		contextutil.GetTraceID(ctx),
		"ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedCR),
		logger.ServiceOperationSuccess,
	)

	return updatedCR, nil
}
