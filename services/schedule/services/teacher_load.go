package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type TeacherLoadService interface {
	platform.BaseService[entities.TeacherLoad]
}

func NewTeacherLoadService(tr repositories.TeacherLoadRepository, eb events.EventBus) TeacherLoadService {
	sc := platform.NewServiceConfig("TeacherLoadService", filepath.Join("data", "teacher_loads.json"), "teacher_load")

	res := &teacherLoadService{
		repository: tr,
	}
	res.BaseService = platform.NewBaseService(sc, tr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.TeacherLoad]{},
	)
	res.logger = res.GetLogger()

	eb.Subscribe(contextutil.SetTraceID(context.Background()), string(events.TeacherLoadRT), res.eventHandler)

	return res
}

type teacherLoadService struct {
	platform.BaseService[entities.TeacherLoad]
	repository repositories.TeacherLoadRepository
	logger     logger.Logger
}

func (s *teacherLoadService) eventHandler(ctx context.Context, b []byte) {
	loadEvent, err := events.FromByteConvertor[events.TeacherLoadRE](b)
	if err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "TeacherLoadRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "TeacherLoadRTHandler",
		fmt.Sprintf("received %s", loadEvent), logger.EventDataReceived,
	)

	trueLoad := entities.TeacherLoad{
		ID:           loadEvent.ID,
		TeacherID:    loadEvent.TeacherID,
		DisciplineID: loadEvent.DisciplineID,
		LessonTypeID: loadEvent.LessonTypeID,
		GroupCount:   loadEvent.GroupCount,
	}

	switch loadEvent.Event {
	case events.EntityCreated:
		s.Add(ctx, trueLoad)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, loadEvent.ID, trueLoad)
	case events.EntityDeleted:
		s.Delete(ctx, loadEvent.ID)
	}
}

func (s *teacherLoadService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	load entities.TeacherLoad,
) (entities.TeacherLoad, error) {
	updatedLoad, err := s.repository.ExternalUpdate(ctx, id, load)
	if err != nil {
		return entities.TeacherLoad{}, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "ExternalUpdate",
			err, logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedLoad), logger.ServiceOperationSuccess)
	return updatedLoad, nil
}
