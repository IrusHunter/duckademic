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
)

type TaskStudentService interface {
	platform.BaseService[entities.TaskStudent]
}

func NewTaskStudentService(
	tsr repositories.TaskStudentRepository,
	tr repositories.TaskRepository,
	sr repositories.StudentRepository,
) TaskStudentService {
	sc := platform.NewServiceConfig(
		"TaskStudentService",
		filepath.Join("data", "task_students.json"),
		"task_student",
	)

	res := &taskStudentService{
		repository:        tsr,
		taskRepository:    tr,
		studentRepository: sr,
	}

	res.BaseService = platform.NewBaseService(sc, tsr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.TaskStudent]{
			platform.OnAddPrepare: res.onAddPrepare,
		},
	)

	res.logger = res.GetLogger()
	return res
}

type taskStudentService struct {
	platform.BaseService[entities.TaskStudent]
	repository        repositories.TaskStudentRepository
	taskRepository    repositories.TaskRepository
	studentRepository repositories.StudentRepository
	logger            logger.Logger
}

func (s *taskStudentService) onAddPrepare(
	ctx context.Context,
	ts *entities.TaskStudent,
) error {
	ts.ID = uuid.New()
	return nil
}

func (s *taskStudentService) Seed(ctx context.Context) error {
	type seedItem struct {
		TaskTitle   string   `json:"task_title"`
		StudentName string   `json:"student_name"`
		Mark        *float64 `json:"mark"`
	}

	var items []seedItem
	if err := jsonutil.ReadFileTo(filepath.Join("data", "task_students.json"), &items); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load task student seed data: %w", err),
			logger.ServiceValidationFailed,
		)
	}

	var lastError error

	for _, item := range items {
		task := s.taskRepository.FindFirstByTitle(ctx, item.TaskTitle)
		if task == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("task %q not found", item.TaskTitle),
				logger.ServiceValidationFailed,
			)
			continue
		}

		student := s.studentRepository.FindFirstByName(ctx, item.StudentName)
		if student == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("student %q not found", item.StudentName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		tsEntity := entities.TaskStudent{
			TaskID:    task.ID,
			StudentID: student.ID,
			Mark:      item.Mark,
		}

		if item.Mark != nil {
			now := time.Now()
			tsEntity.SubmissionTime = &now
		}

		_, err := s.Add(ctx, tsEntity)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to add %s: %w", tsEntity, err),
				logger.ServiceValidationFailed,
			)
			continue
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d task student mappings processed from seed", len(items)),
		logger.ServiceOperationSuccess,
	)

	return lastError
}
