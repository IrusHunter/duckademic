package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

// TeacherSlotPriorityService provides operations for teacher slot priority management.
type TeacherSlotPriorityService interface {
	platform.BaseService[entities.TeacherSlotPriority]
}

// NewTeacherSlotPriorityService creates a new TeacherSlotPriorityService instance.
func NewTeacherSlotPriorityService(
	r repositories.TeacherSlotPriorityRepository,
	tr repositories.TeacherRepository,
	lsr repositories.LessonSlotRepository,
) TeacherSlotPriorityService {
	sc := platform.NewServiceConfig(
		"TeacherSlotPriorityService",
		filepath.Join("data", "teacher_slot_priorities.json"),
		entities.TeacherSlotPriority{}.EntityName(),
	)

	s := &teacherSlotPriorityService{
		repository:           r,
		teacherRepository:    tr,
		lessonSlotRepository: lsr,
	}

	s.BaseService = platform.NewBaseService(sc, r,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.TeacherSlotPriority]{
			platform.OnAddPrepare:   s.onAddPrepare,
			platform.ValidateEntity: s.validateEntity,
		},
	)

	return s
}

type teacherSlotPriorityService struct {
	platform.BaseService[entities.TeacherSlotPriority]
	repository           repositories.TeacherSlotPriorityRepository
	teacherRepository    repositories.TeacherRepository
	lessonSlotRepository repositories.LessonSlotRepository
}

func (s *teacherSlotPriorityService) validateEntity(ctx context.Context, tsp *entities.TeacherSlotPriority) error {
	if err := tsp.ValidatePriority(); err != nil {
		return err
	}

	return nil
}
func (s *teacherSlotPriorityService) onAddPrepare(ctx context.Context, tsp *entities.TeacherSlotPriority) error {
	tsp.ID = uuid.New()

	return nil
}

func (s *teacherSlotPriorityService) Seed(ctx context.Context) error {
	teacherSlotPriorities := []struct {
		TeacherName string                            `json:"teacher_name"`
		Day         int                               `json:"day"`
		Slot        int                               `json:"slot"`
		Priority    entities.TeacherSlotPriorityValue `json:"priority"`
	}{}

	if err := jsonutil.ReadFileTo(
		filepath.Join("data", "teacher_slot_priorities.json"),
		&teacherSlotPriorities,
	); err != nil {
		return s.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"Seed",
			fmt.Errorf("failed to load teacher slot priorities seed data: %w", err),
			logger.ServiceValidationFailed,
		)
	}

	var lastError error

	for _, item := range teacherSlotPriorities {
		teacher := s.teacherRepository.FindByName(ctx, item.TeacherName)
		if teacher == nil {
			lastError = s.GetLogger().LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("teacher %q not found", item.TeacherName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		timeSlot := s.lessonSlotRepository.FindBySlotAndWeekday(ctx, item.Slot, item.Day)
		if timeSlot == nil {
			lastError = s.GetLogger().LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("time slot not found for day=%d slot=%d", item.Day, item.Slot),
				logger.ServiceValidationFailed,
			)
			continue
		}

		truePriority := entities.TeacherSlotPriority{
			TeacherID:  teacher.ID,
			TimeSlotID: timeSlot.ID,
			Priority:   item.Priority,
		}

		_, err := s.Add(ctx, truePriority)
		if err != nil {
			lastError = s.GetLogger().LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("failed to add %+v: %w", truePriority, err),
				logger.ServiceValidationFailed,
			)
			continue
		}
	}

	s.GetLogger().Log(
		contextutil.GetTraceID(ctx),
		"Seed",
		fmt.Sprintf(
			"%d teacher slot priorities processed from seed",
			len(teacherSlotPriorities),
		),
		logger.ServiceOperationSuccess,
	)

	return lastError
}
