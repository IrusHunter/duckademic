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

type SemesterDisciplineService interface {
	platform.BaseService[entities.SemesterDiscipline]
	GetBySemesterIDs(context.Context, []uuid.UUID) ([]entities.Discipline, error)
}

func NewSemesterDisciplineService(
	sr repositories.SemesterDisciplineRepository,
	eb events.EventBus,
) SemesterDisciplineService {
	sc := platform.NewServiceConfig(
		"SemesterDisciplineService",
		filepath.Join("data", "semester_disciplines.json"),
		"semester_discipline",
	)

	res := &semesterDisciplineService{
		repository: sr,
	}

	res.BaseService = platform.NewBaseService(sc, sr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.SemesterDiscipline]{},
	)

	res.logger = res.GetLogger()

	eb.Subscribe(
		contextutil.SetTraceID(context.Background()),
		string(events.SemesterDisciplineRT),
		res.eventHandler,
	)

	return res
}

type semesterDisciplineService struct {
	platform.BaseService[entities.SemesterDiscipline]
	repository repositories.SemesterDisciplineRepository
	logger     logger.Logger
}

func (s *semesterDisciplineService) eventHandler(ctx context.Context, b []byte) {
	sdr, err := events.FromByteConvertor[events.SemesterDisciplineRE](b)
	if err != nil {
		s.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"SemesterDisciplineRTHandler",
			err,
			logger.EventDataReadFailed,
		)
		return
	}

	s.logger.Log(
		contextutil.GetTraceID(ctx),
		"SemesterDisciplineRTHandler",
		fmt.Sprintf("received %s", sdr),
		logger.EventDataReceived,
	)

	trueSD := entities.SemesterDiscipline{
		ID:           sdr.ID,
		SemesterID:   sdr.SemesterID,
		DisciplineID: sdr.DisciplineID,
	}

	switch sdr.Event {
	case events.EntityCreated:
		s.Add(ctx, trueSD)
	case events.EntityDeleted:
		s.Delete(ctx, sdr.ID)
	}
}

func (s *semesterDisciplineService) GetBySemesterIDs(
	ctx context.Context, semesterIDs []uuid.UUID,
) ([]entities.Discipline, error) {
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup
	var mu sync.Mutex

	var result []entities.Discipline
	var lastError error

	for _, semesterID := range semesterIDs {
		wg.Add(1)
		sem <- struct{}{}

		go func(semesterID uuid.UUID) {
			defer wg.Done()
			defer func() { <-sem }()

			disciplines, err := s.repository.GetDisciplinesBySemesterID(ctx, semesterID)
			if err != nil {
				mu.Lock()
				lastError = s.GetLogger().LogAndReturnError(
					contextutil.GetTraceID(ctx),
					"GetBySemesterIDs",
					err,
					logger.ServiceRepositoryFailed,
				)
				mu.Unlock()
				return
			}

			mu.Lock()
			result = append(result, disciplines...)
			mu.Unlock()

		}(semesterID)
	}

	wg.Wait()

	return result, lastError
}
