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

type SemesterDisciplineService interface {
	platform.BaseService[entities.SemesterDiscipline]
}

func NewSemesterDisciplineService(
	sdr repositories.SemesterDisciplineRepository,
	sr repositories.SemesterRepository,
	dr repositories.DisciplineRepository,
	cr repositories.CurriculumRepository,
	eb events.EventBus,
) SemesterDisciplineService {
	sc := platform.NewServiceConfig(
		"SemesterDisciplineService",
		filepath.Join("data", "semester_discipline.json"),
		"semester_discipline",
	)

	res := &semesterDisciplineService{
		repository:           sdr,
		semesterRepository:   sr,
		disciplineRepository: dr,
		curriculumRepository: cr,
	}

	res.BaseService = platform.NewBaseServiceWithEventBus(sc, sdr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.SemesterDiscipline]{
			platform.OnAddPrepare: res.onAddPrepare,
		},
		eb,
	)

	res.logger = res.GetLogger()
	return res
}

type semesterDisciplineService struct {
	platform.BaseService[entities.SemesterDiscipline]
	repository           repositories.SemesterDisciplineRepository
	semesterRepository   repositories.SemesterRepository
	disciplineRepository repositories.DisciplineRepository
	curriculumRepository repositories.CurriculumRepository
	logger               logger.Logger
}

func (s *semesterDisciplineService) onAddPrepare(
	ctx context.Context, semesterDiscipline *entities.SemesterDiscipline,
) error {
	semesterDiscipline.ID = uuid.New()
	return nil
}
func (s *semesterDisciplineService) Seed(ctx context.Context) error {
	type seedItem struct {
		CurriculumName string `json:"curriculum_name"`
		SemesterNumber int    `json:"semester_number"`
		DisciplineName string `json:"discipline_name"`
	}

	var mappings []seedItem
	if err := jsonutil.ReadFileTo(filepath.Join("data", "semester_discipline.json"), &mappings); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load semester discipline seed data: %w", err), logger.ServiceValidationFailed,
		)
	}

	var lastError error

	for _, item := range mappings {
		semesterSlug := fmt.Sprintf("%s-%d", slug.Make(item.CurriculumName), item.SemesterNumber)
		semester := s.semesterRepository.FindBySlug(ctx, semesterSlug)
		if semester == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("semester slug %q not found", semesterSlug), logger.ServiceValidationFailed,
			)
			continue
		}

		discipline := s.disciplineRepository.FindFirstByName(ctx, item.DisciplineName)
		if discipline == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("discipline %q not found", item.DisciplineName), logger.ServiceValidationFailed,
			)
			continue
		}

		sd := entities.SemesterDiscipline{
			SemesterID:   semester.ID,
			DisciplineID: discipline.ID,
		}

		_, err := s.Add(ctx, sd)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to add %s: %w", sd, err), logger.ServiceValidationFailed,
			)
			continue
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d semester discipline mappings processed from seed", len(mappings)), logger.ServiceOperationSuccess,
	)

	return lastError
}

func (s *semesterDisciplineService) Add(
	ctx context.Context, sd entities.SemesterDiscipline,
) (entities.SemesterDiscipline, error) {
	addedSD, err := s.BaseService.Add(ctx, sd)
	if err == nil {
		s.sendChanges(ctx, addedSD, events.EntityCreated)
	}
	return addedSD, err
}

func (s *semesterDisciplineService) Delete(
	ctx context.Context, id uuid.UUID,
) (entities.SemesterDiscipline, error) {
	deletedSD, err := s.BaseService.Delete(ctx, id)
	if err == nil {
		s.sendChanges(ctx, deletedSD, events.EntityDeleted)
	}
	return deletedSD, err
}

func (s *semesterDisciplineService) Update(
	ctx context.Context, id uuid.UUID, sd entities.SemesterDiscipline,
) (entities.SemesterDiscipline, error) {
	updatedSD, err := s.BaseService.Update(ctx, id, sd)
	if err == nil {
		s.sendChanges(ctx, updatedSD, events.EntityUpdated)
	}
	return updatedSD, err
}

func (s *semesterDisciplineService) sendChanges(
	ctx context.Context,
	sd entities.SemesterDiscipline,
	event events.EventType,
) {
	eventSD := events.SemesterDisciplineRE{
		Event:        event,
		ID:           sd.ID,
		SemesterID:   sd.SemesterID,
		DisciplineID: sd.DisciplineID,
	}

	s.BaseService.SendChanges(ctx, eventSD, event, events.SemesterDisciplineRT)
}
