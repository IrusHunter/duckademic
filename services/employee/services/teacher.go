package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/employee/entities"
	"github.com/IrusHunter/duckademic/services/employee/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type TeacherService interface {
	platform.BaseService[entities.Teacher]
}

func NewTeacherService(
	tr repositories.TeacherRepository,
	arr repositories.AcademicRankRepository,
	adr repositories.AcademicDegreeRepository,
	er repositories.EmployeeRepository,
	eb events.EventBus,
) TeacherService {
	sc := platform.NewServiceConfig(
		"TeacherService",
		filepath.Join("data", "teachers.json"),
		"teacher",
	)

	res := &teacherService{
		repository:               tr,
		academicRankRepository:   arr,
		academicDegreeRepository: adr,
		employeeRepository:       er,
		eventBus:                 eb,
	}

	res.BaseService = platform.NewBaseServiceWithEventBus(sc, tr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Teacher]{
			platform.OnAddPrepare:   res.onAddPrepare,
			platform.ValidateEntity: res.validateEntity,
		},
		eb,
	)

	res.logger = res.GetLogger()

	return res
}

type teacherService struct {
	platform.BaseService[entities.Teacher]
	repository               repositories.TeacherRepository
	academicRankRepository   repositories.AcademicRankRepository
	academicDegreeRepository repositories.AcademicDegreeRepository
	employeeRepository       repositories.EmployeeRepository
	logger                   logger.Logger
	eventBus                 events.EventBus
}

func (s *teacherService) validateEntity(ctx context.Context, teacher *entities.Teacher) error {
	if err := teacher.ValidateEmail(); err != nil {
		return err
	}

	return nil
}
func (s *teacherService) onAddPrepare(ctx context.Context, teacher *entities.Teacher) error {
	return nil
}
func (s *teacherService) hardDeleteCheck(ctx context.Context, teacher *entities.Teacher) error {
	return fmt.Errorf("plug")
}

type seedTeacher struct {
	FirstName           string  `json:"first_name"`
	LastName            string  `json:"last_name"`
	Email               string  `json:"email"`
	AcademicRankTitle   string  `json:"academic_rank_title"`
	AcademicDegreeTitle *string `json:"academic_rank_degree,omitempty"`
}

func (s *teacherService) Seed(ctx context.Context) error {
	teachers := []seedTeacher{}
	if err := jsonutil.ReadFileTo(filepath.Join("data", "teachers.json"), &teachers); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load teachers seed data: %w", err), logger.ServiceValidationFailed,
		)
	}

	var lastError error
	for _, teacher := range teachers {
		trueTeacher := entities.Teacher{Email: teacher.Email}
		employee := s.employeeRepository.FindFirstByName(ctx, teacher.FirstName, teacher.LastName)
		if employee == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("employee %s %s not found", teacher.FirstName, teacher.LastName), logger.ServiceValidationFailed,
			)
			continue
		}
		trueTeacher.EmployeeID = employee.ID

		academicRank := s.academicRankRepository.FindByTitle(ctx, teacher.AcademicRankTitle)
		if academicRank == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("academic rank with title %q not found", teacher.AcademicRankTitle), logger.ServiceValidationFailed,
			)
			continue
		}
		trueTeacher.AcademicRankID = academicRank.ID

		if teacher.AcademicDegreeTitle != nil {
			academicDegree := s.academicDegreeRepository.FindByTitle(ctx, *teacher.AcademicDegreeTitle)
			if academicDegree == nil {
				lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
					fmt.Errorf("academic degree with title %q not found", *teacher.AcademicDegreeTitle),
					logger.ServiceValidationFailed,
				)
				continue
			}
			trueTeacher.AcademicDegreeID = &academicDegree.ID
		}

		_, err := s.Add(ctx, trueTeacher)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to add %s: %w", trueTeacher.String(), err), logger.ServiceValidationFailed,
			)
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d teachers added successfully", len(teachers)), logger.ServiceOperationSuccess,
	)
	return lastError
}

func (s *teacherService) Add(
	ctx context.Context, teacher entities.Teacher,
) (entities.Teacher, error) {
	addedT, err := s.BaseService.Add(ctx, teacher)
	if err == nil {
		s.sendChanges(ctx, addedT, events.EntityCreated)
	}
	return addedT, err
}
func (s *teacherService) Delete(
	ctx context.Context, id uuid.UUID,
) (entities.Teacher, error) {
	deletedT, err := s.BaseService.Delete(ctx, id)
	if err == nil {
		s.sendChanges(ctx, deletedT, events.EntityDeleted)
	}
	return deletedT, err
}
func (s *teacherService) Update(
	ctx context.Context, id uuid.UUID, teacher entities.Teacher,
) (entities.Teacher, error) {
	updatedT, err := s.BaseService.Update(ctx, id, teacher)
	if err == nil {
		s.sendChanges(ctx, updatedT, events.EntityUpdated)
	}
	return updatedT, err
}
func (s *teacherService) sendChanges(ctx context.Context, teacher entities.Teacher, event events.EventType) {
	filledT := s.repository.Fill(ctx, teacher.EmployeeID)

	if filledT == nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SendChanges",
			fmt.Errorf("failed to fill %s: not found", teacher), logger.ServiceDataFetchFailed,
		)
		return
	}

	eventT := events.TeacherRE{
		Event:          events.EntityCreated,
		ID:             filledT.EmployeeID,
		Slug:           filledT.Employee.Slug,
		Name:           filledT.Employee.GetShortFullName(),
		AcademicRankID: filledT.AcademicRankID,
	}

	s.BaseService.SendChanges(ctx, eventT, event, events.TeacherRT)
}
