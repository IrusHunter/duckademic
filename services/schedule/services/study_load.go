package services

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type StudyLoadService interface {
	platform.BaseService[entities.StudyLoad]
	AddMultiple(context.Context, []entities.StudyLoad) error
}

func NewStudyLoadService(
	sr repositories.StudyLoadRepository,
) StudyLoadService {
	sc := platform.NewServiceConfig(
		"StudyLoadService",
		filepath.Join("data", "study_loads.json"),
		"study load",
	)

	res := &studyLoadService{
		repository: sr,
	}
	res.BaseService = platform.NewBaseService(
		sc,
		sr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.StudyLoad]{},
	)

	return res
}

type studyLoadService struct {
	platform.BaseService[entities.StudyLoad]
	repository repositories.StudyLoadRepository
}

func (s *studyLoadService) AddMultiple(ctx context.Context, studyLoads []entities.StudyLoad) error {
	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var lastError error

	for i, studyLoad := range studyLoads {
		wg.Add(1)
		sem <- struct{}{}

		go func(i int, studyLoad entities.StudyLoad) {
			defer wg.Done()
			defer func() { <-sem }()

			_, err := s.Add(ctx, studyLoad)
			if err != nil {
				mu.Lock()
				lastError = s.GetLogger().LogAndReturnError(contextutil.GetTraceID(ctx), "AddMultiple",
					fmt.Errorf("failed to insert at index [%d]: %w", i, err), logger.ServiceValidationFailed)
				mu.Unlock()
			}
		}(i, studyLoad)

	}

	wg.Wait()
	return lastError
}
