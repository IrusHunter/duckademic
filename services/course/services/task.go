package services

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/services/course/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type TaskService interface {
	platform.BaseService[entities.Task]
}

func NewTaskService(
	tr repositories.TaskRepository,
	cr repositories.CourseRepository,
) TaskService {
	sc := platform.NewServiceConfig(
		"TaskService",
		filepath.Join("data", "tasks.json"),
		"task",
	)

	res := &taskService{
		repository: tr,
		courseRepo: cr,
	}

	res.BaseService = platform.NewBaseService(sc, tr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Task]{
			platform.OnAddPrepare: res.onAddPrepare,
		},
	)

	res.logger = res.GetLogger()
	return res
}

type taskService struct {
	platform.BaseService[entities.Task]
	repository repositories.TaskRepository
	courseRepo repositories.CourseRepository
	logger     logger.Logger
}

func (s *taskService) onAddPrepare(ctx context.Context, t *entities.Task) error {
	t.Slug = slug.Make(t.Title)
	t.ID = uuid.New()
	return nil
}

func (s *taskService) Seed(ctx context.Context) error {
	type seedItem struct {
		CourseName  string    `json:"course_name"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		MaxMark     float64   `json:"max_mark"`
		Deadline    time.Time `json:"deadline"`
	}

	var items []seedItem
	if err := jsonutil.ReadFileTo(filepath.Join("data", "tasks.json"), &items); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load task seed data: %w", err),
			logger.ServiceValidationFailed,
		)
	}

	var lastError error

	for _, item := range items {
		course := s.courseRepo.FindFirstByName(ctx, item.CourseName)
		if course == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("course %q not found", item.CourseName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		task := entities.Task{
			CourseID:    course.ID,
			Title:       item.Title,
			Description: item.Description,
			MaxMark:     item.MaxMark,
			Deadline:    item.Deadline,
		}

		_, err := s.Add(ctx, task)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to add %s: %w", task, err),
				logger.ServiceValidationFailed,
			)
			continue
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d tasks processed from seed", len(items)),
		logger.ServiceOperationSuccess,
	)

	return lastError
}
