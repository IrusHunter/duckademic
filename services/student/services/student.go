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

type StudentService interface {
	platform.BaseService[entities.Student]
}

func NewStudentService(
	sr repositories.StudentRepository,
	eb events.EventBus,
) StudentService {
	sc := platform.NewServiceConfig(
		"StudentService",
		filepath.Join("data", "students.json"),
		"student",
	)

	res := &studentService{
		repository: sr,
		eventBus:   eb,
	}

	res.BaseService = platform.NewBaseServiceWithEventBus(sc, sr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Student]{
			platform.OnAddPrepare:    res.onAddPrepare,
			platform.ValidateEntity:  res.validateEntity,
			platform.HardDeleteCheck: res.hardDeleteCheck,
		},
		eb,
	)

	res.logger = res.GetLogger()

	return res
}

type studentService struct {
	platform.BaseService[entities.Student]
	repository repositories.StudentRepository
	logger     logger.Logger
	eventBus   events.EventBus
}

func (s *studentService) validateEntity(ctx context.Context, student *entities.Student) error {
	if err := student.ValidateFirstName(); err != nil {
		return err
	}
	if err := student.ValidateLastName(); err != nil {
		return err
	}
	if err := student.ValidateEmail(); err != nil {
		return err
	}
	return nil
}

func (s *studentService) onAddPrepare(ctx context.Context, student *entities.Student) error {
	slug := slug.Make(student.GetFullName())
	if other := s.repository.FindBySlug(ctx, slug); other != nil {
		return fmt.Errorf("employee with slug %q already exists", slug)
	}
	student.ID = uuid.New()
	student.Slug = slug

	return nil
}

func (s *studentService) hardDeleteCheck(ctx context.Context, student *entities.Student) error {
	return fmt.Errorf("plug")
}

func (s *studentService) Seed(ctx context.Context) error {
	students := []entities.Student{}

	if err := jsonutil.ReadFileTo(filepath.Join("data", "students.json"), &students); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load students seed data: %w", err),
			logger.ServiceValidationFailed,
		)
	}

	var lastError error
	for _, student := range students {
		_, err := s.Add(ctx, student)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to add %s: %w", student.String(), err),
				logger.ServiceValidationFailed,
			)
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d students added successfully", len(students)),
		logger.ServiceOperationSuccess,
	)

	return lastError
}

func (s *studentService) Add(
	ctx context.Context, student entities.Student,
) (entities.Student, error) {
	added, err := s.BaseService.Add(ctx, student)
	if err == nil {
		s.sendChanges(ctx, added, events.EntityCreated)
	}
	return added, err
}
func (s *studentService) Delete(
	ctx context.Context, id uuid.UUID,
) (entities.Student, error) {
	deleted, err := s.BaseService.Delete(ctx, id)
	if err == nil {
		s.sendChanges(ctx, deleted, events.EntityDeleted)
	}
	return deleted, err
}
func (s *studentService) Update(
	ctx context.Context, id uuid.UUID, student entities.Student,
) (entities.Student, error) {
	updated, err := s.BaseService.Update(ctx, id, student)
	if err == nil {
		s.sendChanges(ctx, updated, events.EntityUpdated)
	}
	return updated, err
}

func (s *studentService) sendChanges(
	ctx context.Context,
	student entities.Student,
	eventType events.EventType,
) {
	eventS := events.StudentRE{
		Event: eventType,
		ID:    student.ID,
		Name:  student.GetShortFullName(),
	}

	s.BaseService.SendChanges(ctx, eventS, eventType, events.StudentRT)
}
