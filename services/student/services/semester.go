package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/student/entities"
	"github.com/IrusHunter/duckademic/services/student/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type SemesterService interface {
	platform.BaseService[entities.Semester]
}

func NewSemesterService(
	sr repositories.SemesterRepository,
	eb events.EventBus,
) SemesterService {
	sc := platform.NewServiceConfig(
		"SemesterService",
		filepath.Join("data", "semesters.json"),
		"semester",
	)

	res := &semesterService{
		repository: sr,
	}

	res.BaseService = platform.NewBaseService(sc, sr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Semester]{},
	)

	res.logger = res.GetLogger()

	eb.Subscribe(contextutil.SetTraceID(context.Background()), string(events.SemesterRT), res.eventHandler)

	return res
}

type semesterService struct {
	platform.BaseService[entities.Semester]
	repository repositories.SemesterRepository
	logger     logger.Logger
}

func (s *semesterService) eventHandler(ctx context.Context, b []byte) {
	sr, err := events.FromByteConvertor[events.SemesterRE](b)
	if err != nil {
		s.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"SemesterRTHandler",
			err,
			logger.EventDataReadFailed,
		)
		return
	}

	s.logger.Log(
		contextutil.GetTraceID(ctx),
		"SemesterRTHandler",
		fmt.Sprintf("received %s", sr),
		logger.EventDataReceived,
	)

	trueSR := entities.Semester{
		ID:           sr.ID,
		Slug:         sr.Slug,
		CurriculumID: sr.CurriculumID,
		Number:       sr.Number,
	}

	switch sr.Event {
	case events.EntityCreated:
		s.Add(ctx, trueSR)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, sr.ID, trueSR)
	case events.EntityDeleted:
		s.Delete(ctx, sr.ID)
	}
}

func (s *semesterService) Seed(ctx context.Context) error {
	semesters := []struct {
		CurriculumName string `json:"curriculum_name"`
		Number         int    `json:"number"`
	}{}

	if err := jsonutil.ReadFileTo(filepath.Join("data", "semesters.json"), &semesters); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load semesters seed data: %w", err), logger.ServiceValidationFailed,
		)
	}

	var lastError error
	for _, semester := range semesters {
		semesterSlug := fmt.Sprintf("%s-%d", slug.Make(semester.CurriculumName), semester.Number)
		semester := s.repository.FindBySlug(ctx, semesterSlug)
		if semester == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("semester slug %q not found", semesterSlug), logger.ServiceValidationFailed,
			)
			continue
		}

		_, err := s.Update(ctx, semester.ID, *semester)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to update %s: %w", semester.String(), err),
				logger.ServiceValidationFailed,
			)
			continue
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d semesters processed from seed", len(semesters)), logger.ServiceOperationSuccess,
	)

	return lastError
}
func (s *semesterService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	semester entities.Semester,
) (entities.Semester, error) {
	updatedS, err := s.repository.ExternalUpdate(ctx, id, semester)
	if err != nil {
		return entities.Semester{}, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "ExternalUpdate",
			err, logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedS), logger.ServiceOperationSuccess)
	return updatedS, nil
}
