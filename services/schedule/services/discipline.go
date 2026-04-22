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

type DisciplineService interface {
	platform.BaseService[entities.Discipline]
	ToGeneratorDisciplines(context.Context, []entities.Discipline) []GeneratorDiscipline
	ExtractIDs([]entities.Discipline) []uuid.UUID
}

func NewDisciplineService(dr repositories.DisciplineRepository, eb events.EventBus) DisciplineService {
	sc := platform.NewServiceConfig("DisciplineService", filepath.Join("data", "disciplines.json"), "discipline")

	res := &disciplineService{
		repository: dr,
	}
	res.BaseService = platform.NewBaseService(sc, dr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Discipline]{},
	)
	res.logger = res.GetLogger()

	eb.Subscribe(contextutil.SetTraceID(context.Background()), string(events.DisciplineRT), res.eventHandler)

	return res
}

type disciplineService struct {
	platform.BaseService[entities.Discipline]
	repository repositories.DisciplineRepository
	logger     logger.Logger
}

func (s *disciplineService) eventHandler(ctx context.Context, b []byte) {
	dr, err := events.FromByteConvertor[events.DisciplineRE](b)
	if err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "DisciplineRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "DisciplineRTHandler",
		fmt.Sprintf("received %s", dr), logger.EventDataReceived,
	)

	trueDR := entities.Discipline{
		ID:   dr.ID,
		Name: dr.Name,
		Slug: dr.Slug,
	}

	switch dr.Event {
	case events.EntityCreated:
		s.Add(ctx, trueDR)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, dr.ID, trueDR)
	case events.EntityDeleted:
		s.Delete(ctx, dr.ID)
	}
}

func (s *disciplineService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	discipline entities.Discipline,
) (entities.Discipline, error) {
	updatedDR, err := s.repository.ExternalUpdate(ctx, id, discipline)
	if err != nil {
		return entities.Discipline{}, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "ExternalUpdate",
			err, logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedDR), logger.ServiceOperationSuccess)
	return updatedDR, nil
}

type GeneratorDiscipline struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (s *disciplineService) ToGeneratorDisciplines(ctx context.Context, d []entities.Discipline) []GeneratorDiscipline {
	res := make([]GeneratorDiscipline, 0, len(d))

	for _, discipline := range d {
		res = append(res, GeneratorDiscipline{
			ID:   discipline.ID,
			Name: discipline.Name,
		})
	}

	return res
}
func (s *disciplineService) ExtractIDs(d []entities.Discipline) []uuid.UUID {
	res := make([]uuid.UUID, 0, len(d))

	for _, discipline := range d {
		res = append(res, discipline.ID)
	}

	return res
}
