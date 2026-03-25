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
	semR repositories.SemesterRepository,
	eb events.EventBus,
) StudentService {
	sc := platform.NewServiceConfig(
		"StudentService",
		filepath.Join("data", "students.json"),
		"student",
	)

	res := &studentService{
		repository:         sr,
		semesterRepository: semR,
		eventBus:           eb,
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
	repository         repositories.StudentRepository
	semesterRepository repositories.SemesterRepository
	logger             logger.Logger
	eventBus           events.EventBus
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
	studentsData := []struct {
		FirstName      string  `json:"first_name"`
		LastName       string  `json:"last_name"`
		MiddleName     *string `json:"middle_name,omitempty"`
		Email          string  `json:"email"`
		PhoneNumber    *string `json:"phone_number,omitempty"`
		CurriculumName string  `json:"curriculum_name"`
		SemesterNumber int     `json:"semester_number"`
	}{}

	if err := jsonutil.ReadFileTo(filepath.Join("data", "students.json"), &studentsData); err != nil {
		return s.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"Seed",
			fmt.Errorf("failed to load students seed data: %w", err),
			logger.ServiceValidationFailed,
		)
	}

	var lastError error
	for _, item := range studentsData {
		semesterSlug := fmt.Sprintf("%s-%d", slug.Make(item.CurriculumName), item.SemesterNumber)
		semester := s.semesterRepository.FindBySlug(ctx, semesterSlug)
		if semester == nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("semester slug %q not found", semesterSlug),
				logger.ServiceValidationFailed,
			)
			continue
		}

		student := entities.Student{
			FirstName:   item.FirstName,
			LastName:    item.LastName,
			MiddleName:  item.MiddleName,
			Email:       item.Email,
			PhoneNumber: item.PhoneNumber,
			SemesterID:  semester.ID,
		}

		_, err := s.Add(ctx, student)
		if err != nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("failed to add %s %s: %w", student.FirstName, student.LastName, err),
				logger.ServiceValidationFailed,
			)
		}
	}

	s.logger.Log(
		contextutil.GetTraceID(ctx),
		"Seed",
		fmt.Sprintf("%d students processed", len(studentsData)),
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
		Slug:  student.Slug,
		Name:  student.GetShortFullName(),
	}

	s.BaseService.SendChanges(ctx, eventS, eventType, events.StudentRT)
}
