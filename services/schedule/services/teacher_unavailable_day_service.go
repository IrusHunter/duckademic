package services

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

// TeacherUnavailableDayService provides operations for teacher unavailable day management.
type TeacherUnavailableDayService interface {
	platform.BaseService[entities.TeacherUnavailableDay]
}

// NewTeacherUnavailableDayService creates a new TeacherUnavailableDayService instance.
func NewTeacherUnavailableDayService(
	r repositories.TeacherUnavailableDayRepository,
	tr repositories.TeacherRepository,
) TeacherUnavailableDayService {
	sc := platform.NewServiceConfig(
		"TeacherUnavailableDayService",
		filepath.Join("data", "teacher_unavailable_days.json"),
		entities.TeacherUnavailableDay{}.EntityName(),
	)

	s := &teacherUnavailableDayService{
		repository:        r,
		teacherRepository: tr,
	}

	s.BaseService = platform.NewBaseService(sc, r,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.TeacherUnavailableDay]{
			platform.OnAddPrepare: s.onAddPrepare,
		},
	)

	return s
}

type teacherUnavailableDayService struct {
	platform.BaseService[entities.TeacherUnavailableDay]
	repository        repositories.TeacherUnavailableDayRepository
	teacherRepository repositories.TeacherRepository
}

func (s *teacherUnavailableDayService) validateEntity(ctx context.Context, tud *entities.TeacherUnavailableDay) error {
	if err := tud.ValidateDay(); err != nil {
		return err
	}

	return nil
}
func (s *teacherUnavailableDayService) onAddPrepare(ctx context.Context, tud *entities.TeacherUnavailableDay) error {
	tud.ID = uuid.New()

	return nil
}

func (s *teacherUnavailableDayService) Seed(ctx context.Context) error {
	teacherUnavailableDays := []struct {
		TeacherName string `json:"teacher_name"`
		Day         string `json:"day"`
	}{}

	if err := jsonutil.ReadFileTo(
		filepath.Join("data", "teacher_unavailable_days.json"),
		&teacherUnavailableDays,
	); err != nil {
		return s.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"Seed",
			fmt.Errorf("failed to load teacher unavailable day seed data: %w", err),
			logger.ServiceValidationFailed,
		)
	}

	var lastError error

	for _, item := range teacherUnavailableDays {
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

		day, err := time.Parse(time.DateOnly, item.Day)
		if err != nil {
			lastError = s.GetLogger().LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("invalid day format %q: %w", item.Day, err),
				logger.ServiceValidationFailed,
			)
			continue
		}

		trueUnavailableDay := entities.TeacherUnavailableDay{
			TeacherID: teacher.ID,
			Day:       day,
		}

		_, err = s.Add(ctx, trueUnavailableDay)
		if err != nil {
			lastError = s.GetLogger().LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("failed to add %+v: %w", trueUnavailableDay, err),
				logger.ServiceValidationFailed,
			)
			continue
		}
	}

	s.GetLogger().Log(
		contextutil.GetTraceID(ctx),
		"Seed",
		fmt.Sprintf(
			"%d teacher unavailable days processed from seed",
			len(teacherUnavailableDays),
		),
		logger.ServiceOperationSuccess,
	)

	return lastError
}
