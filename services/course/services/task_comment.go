package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/services/course/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type TaskCommentService interface {
	platform.BaseService[entities.TaskComment]
	GetForTask(ctx context.Context, taskID uuid.UUID, isPrivate bool, studentID *uuid.UUID) ([]entities.TaskComment, error)
	AddComment(ctx context.Context, taskID, authorID uuid.UUID, role, body string, isPrivate bool, studentID *uuid.UUID) (entities.TaskComment, error)
}

func NewTaskCommentService(
	repo repositories.TaskCommentRepository,
	studentRepo repositories.StudentRepository,
	teacherRepo repositories.TeacherRepository,
) TaskCommentService {
	sc := platform.NewServiceConfig(
		"TaskCommentService",
		filepath.Join("data", "task_comments.json"),
		"task comment",
	)

	svc := &taskCommentService{
		repository:        repo,
		studentRepository: studentRepo,
		teacherRepository: teacherRepo,
	}
	svc.BaseService = platform.NewBaseService(sc, repo,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.TaskComment]{
			platform.OnAddPrepare: svc.onAddPrepare,
		},
	)
	return svc
}

type taskCommentService struct {
	platform.BaseService[entities.TaskComment]
	repository        repositories.TaskCommentRepository
	studentRepository repositories.StudentRepository
	teacherRepository repositories.TeacherRepository
}

func (s *taskCommentService) onAddPrepare(_ context.Context, c *entities.TaskComment) error {
	c.ID = uuid.New()
	return nil
}

func (s *taskCommentService) GetForTask(ctx context.Context, taskID uuid.UUID, isPrivate bool, studentID *uuid.UUID) ([]entities.TaskComment, error) {
	return s.repository.GetForTask(ctx, taskID, isPrivate, studentID)
}

func (s *taskCommentService) AddComment(
	ctx context.Context,
	taskID, authorID uuid.UUID,
	role, body string,
	isPrivate bool,
	studentID *uuid.UUID,
) (entities.TaskComment, error) {
	if body == "" {
		return entities.TaskComment{}, s.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx), "AddComment",
			fmt.Errorf("comment body cannot be empty"),
			logger.ServiceValidationFailed,
		)
	}

	authorName, err := s.resolveAuthorName(ctx, authorID, role)
	if err != nil {
		return entities.TaskComment{}, err
	}

	comment := entities.TaskComment{
		TaskID:     taskID,
		AuthorID:   authorID,
		AuthorName: authorName,
		Body:       body,
		IsPrivate:  isPrivate,
		StudentID:  studentID,
	}

	return s.Add(ctx, comment)
}

func (s *taskCommentService) resolveAuthorName(ctx context.Context, authorID uuid.UUID, role string) (string, error) {
	if role == "teacher" {
		teacher := s.teacherRepository.FindByID(ctx, authorID)
		if teacher == nil {
			return "", s.GetLogger().LogAndReturnError(
				contextutil.GetTraceID(ctx), "resolveAuthorName",
				fmt.Errorf("teacher %s not found", authorID),
				logger.ServiceDataFetchFailed,
			)
		}
		return teacher.Name, nil
	}

	student := s.studentRepository.FindByID(ctx, authorID)
	if student == nil {
		return "", s.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx), "resolveAuthorName",
			fmt.Errorf("student %s not found", authorID),
			logger.ServiceDataFetchFailed,
		)
	}
	return student.Name, nil
}
