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

// AcademicRankService provides operations to initialize and manage academic ranks.
type AcademicRankService interface {
	platform.BaseService[entities.AcademicRank]
}

// NewAcademicRankService creates a new AcademicRankService instance.
//
// It requires an academic rank repository (arr) and an event bus (eb).
func NewAcademicRankService(arr repositories.AcademicRankRepository, eb events.EventBus) AcademicRankService {
	sc := platform.NewServiceConfig("AcademicRankService", filepath.Join("data", "academic_ranks.json"), "academic_rank")

	res := &academicRankService{
		repository: arr,
	}
	res.BaseService = platform.NewBaseService(sc, arr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.AcademicRank]{
			platform.ValidateEntity: res.validateEntity,
		},
	)
	res.logger = res.GetLogger()

	eb.Subscribe(contextutil.SetTraceID(context.Background()), string(events.AcademicRankRT), res.eventHandler)

	return res
}

type academicRankService struct {
	platform.BaseService[entities.AcademicRank]
	repository repositories.AcademicRankRepository
	logger     logger.Logger
}

func (s *academicRankService) validateEntity(ctx context.Context, academicRank *entities.AcademicRank) error {
	if err := academicRank.ValidateTitle(); err != nil {
		return err
	}

	return nil
}
func (s *academicRankService) eventHandler(ctx context.Context, b []byte) {
	ar, err := events.FromByteConvertor[events.AcademicRankRE](b)
	if err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "AcademicRankRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "AcademicRankRTHandler",
		fmt.Sprintf("received %s", ar), logger.EventDataReceived,
	)

	trueAR := entities.AcademicRank{
		ID:    ar.ID,
		Title: ar.Title,
		Slug:  ar.Slug,
	}

	switch ar.Event {
	case events.EntityCreated:
		s.Add(ctx, trueAR)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, ar.ID, trueAR)
	case events.EntityDeleted:
		s.Delete(ctx, ar.ID)
	}
}

func (s *academicRankService) Seed(ctx context.Context) error {
	academicRanks := []entities.AcademicRank{}
	if err := jsonutil.ReadFileTo(filepath.Join("data", "academic_ranks.json"), &academicRanks); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load academic ranks seed data: %w", err), logger.ServiceDataFetchFailed,
		)
	}

	var lastError error
	for _, academicRank := range academicRanks {
		tar := s.repository.FindByTitle(ctx, academicRank.Title)
		if tar == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("academic rank with title %q not found", academicRank.Title), logger.ServiceDataFetchFailed,
			)
			continue
		}

		_, err := s.Update(ctx, tar.ID, academicRank)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to update %s: %w", academicRank, err), logger.ServiceRepositoryFailed,
			)
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d academic ranks updates successfully", len(academicRanks)), logger.ServiceOperationSuccess,
	)
	return lastError
}
func (s *academicRankService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	academicRank entities.AcademicRank,
) (entities.AcademicRank, error) {
	updatedAR, err := s.repository.ExternalUpdate(ctx, id, academicRank)
	if err != nil {
		return entities.AcademicRank{}, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "ExternalUpdate",
			err, logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedAR), logger.ServiceOperationSuccess)
	return updatedAR, nil
}
