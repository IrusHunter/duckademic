package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/services/course/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type TeacherService interface {
	platform.BaseService[entities.Teacher]
}

func NewTeacherService(tr repositories.TeacherRepository, eb events.EventBus) TeacherService {
	sc := platform.NewServiceConfig("TeacherService", filepath.Join("data", "teachers.json"), entities.Teacher{}.EntityName())

	res := &teacherService{
		repository: tr,
	}
	res.BaseService = platform.NewBaseService(sc, tr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Teacher]{},
	)
	res.logger = res.GetLogger()

	eb.Subscribe(contextutil.SetTraceID(context.Background()), string(events.TeacherRT), res.eventHandler)

	return res
}

type teacherService struct {
	platform.BaseService[entities.Teacher]
	repository repositories.TeacherRepository
	logger     logger.Logger
}

func (s *teacherService) eventHandler(ctx context.Context, b []byte) {
	t, err := events.FromByteConvertor[events.TeacherRE](b)
	if err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "TeacherRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "TeacherRTHandler",
		fmt.Sprintf("received %s", t), logger.EventDataReceived,
	)

	trueT := entities.Teacher{
		ID:   t.ID,
		Slug: t.Slug,
		Name: t.Name,
	}

	switch t.Event {
	case events.EntityCreated:
		s.Add(ctx, trueT)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, t.ID, trueT)
	case events.EntityDeleted:
		s.Delete(ctx, t.ID)
	}
}

func (s *teacherService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	teacher entities.Teacher,
) (entities.Teacher, error) {
	updatedT, err := s.repository.ExternalUpdate(ctx, id, teacher)
	if err != nil {
		return entities.Teacher{}, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "ExternalUpdate",
			err, logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedT), logger.ServiceOperationSuccess)
	return updatedT, nil
}
