package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/teacher_load/entities"
	"github.com/IrusHunter/duckademic/services/teacher_load/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type TeacherLoadService interface {
	platform.BaseService[entities.TeacherLoad]
}

func NewTeacherLoadService(
	tlr repositories.TeacherLoadRepository,
	tr repositories.TeacherRepository,
	dr repositories.DisciplineRepository,
	ltr repositories.LessonTypeRepository,
	gcr repositories.GroupCohortRepository,
	eb events.EventBus,
) TeacherLoadService {

	tc := platform.NewServiceConfig(
		"TeacherLoadService",
		filepath.Join("data", "teacher_loads.json"),
		"teacher_load",
	)

	res := &teacherLoadService{
		repository:            tlr,
		teacherRepository:     tr,
		disciplineRepository:  dr,
		lessonTypeRepository:  ltr,
		groupCohortRepository: gcr,
		eventBus:              eb,
	}

	res.BaseService = platform.NewBaseServiceWithEventBus(tc, tlr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.TeacherLoad]{
			platform.OnAddPrepare:   res.onAddPrepare,
			platform.ValidateEntity: res.validateEntity,
		},
		eb,
	)

	res.logger = res.GetLogger()

	return res
}

type teacherLoadService struct {
	platform.BaseService[entities.TeacherLoad]
	repository            repositories.TeacherLoadRepository
	teacherRepository     repositories.TeacherRepository
	disciplineRepository  repositories.DisciplineRepository
	lessonTypeRepository  repositories.LessonTypeRepository
	groupCohortRepository repositories.GroupCohortRepository
	logger                logger.Logger
	eventBus              events.EventBus
}

func (s *teacherLoadService) validateEntity(ctx context.Context, tl *entities.TeacherLoad) error {
	if err := tl.ValidateGroupCount(); err != nil {
		return nil
	}

	return nil
}

func (s *teacherLoadService) onAddPrepare(ctx context.Context, tl *entities.TeacherLoad) error {
	tl.ID = uuid.New()
	return nil
}

func (s *teacherLoadService) Seed(ctx context.Context) error {
	teacherLoads := []struct {
		TeacherName     string `json:"teacher_name"`
		DisciplineName  string `json:"discipline_name"`
		LessonTypeName  string `json:"lesson_type_name"`
		GroupCohortName string `json:"group_cohort_name"`
		GroupCount      int    `json:"group_count"`
	}{}

	if err := jsonutil.ReadFileTo(filepath.Join("data", "teacher_loads.json"), &teacherLoads); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load teacher loads seed data: %w", err),
			logger.ServiceValidationFailed,
		)
	}

	var lastError error
	for _, item := range teacherLoads {
		teacher := s.teacherRepository.FindByName(ctx, item.TeacherName)
		if teacher == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("teacher %q not found", item.TeacherName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		discipline := s.disciplineRepository.FindFirstByName(ctx, item.DisciplineName)
		if discipline == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("discipline %q not found", item.DisciplineName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		lessonType := s.lessonTypeRepository.FindFirstByName(ctx, item.LessonTypeName)
		if lessonType == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("lesson type %q not found", item.LessonTypeName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		groupCohort := s.groupCohortRepository.FindFirstByName(ctx, item.GroupCohortName)
		if groupCohort == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("group cohort %q not found", item.GroupCohortName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		trueTL := entities.TeacherLoad{
			TeacherID:     teacher.ID,
			DisciplineID:  discipline.ID,
			LessonTypeID:  lessonType.ID,
			GroupCohortID: groupCohort.ID,
			GroupCount:    item.GroupCount,
		}

		_, err := s.Add(ctx, trueTL)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to add %s: %w", trueTL, err),
				logger.ServiceValidationFailed,
			)
			continue
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d teacher loads processed from seed", len(teacherLoads)),
		logger.ServiceOperationSuccess,
	)

	return lastError
}

func (s *teacherLoadService) Add(
	ctx context.Context, tl entities.TeacherLoad,
) (entities.TeacherLoad, error) {
	addedTL, err := s.BaseService.Add(ctx, tl)
	if err == nil {
		s.sendChanges(ctx, addedTL, events.EntityCreated)
	}
	return addedTL, err
}
func (s *teacherLoadService) Delete(
	ctx context.Context, id uuid.UUID,
) (entities.TeacherLoad, error) {
	deletedTL, err := s.BaseService.Delete(ctx, id)
	if err == nil {
		s.sendChanges(ctx, deletedTL, events.EntityDeleted)
	}
	return deletedTL, err
}
func (s *teacherLoadService) Update(
	ctx context.Context, id uuid.UUID, tl entities.TeacherLoad,
) (entities.TeacherLoad, error) {
	updatedTL, err := s.BaseService.Update(ctx, id, tl)
	if err == nil {
		s.sendChanges(ctx, updatedTL, events.EntityUpdated)
	}
	return updatedTL, err
}

func (s *teacherLoadService) sendChanges(
	ctx context.Context,
	tl entities.TeacherLoad,
	event events.EventType,
) {
	eventTL := events.TeacherLoadRE{
		Event:         event,
		ID:            tl.ID,
		TeacherID:     tl.TeacherID,
		DisciplineID:  tl.DisciplineID,
		LessonTypeID:  tl.LessonTypeID,
		GroupCohortID: tl.GroupCohortID,
		GroupCount:    tl.GroupCount,
	}

	s.BaseService.SendChanges(ctx, eventTL, event, events.TeacherLoadRT)
}
