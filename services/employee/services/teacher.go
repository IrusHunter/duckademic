package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/employees/entities"
	"github.com/IrusHunter/duckademic/services/employees/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type TeacherService interface {
	platform.BaseService[entities.Teacher]
}

func NewTeacherService(
	tr repositories.TeacherRepository,
	arr repositories.AcademicRankRepository,
	adr repositories.AcademicDegreeRepository,
	er repositories.EmployeeRepository,
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
	}

	res.BaseService = platform.NewBaseService(
		sc,
		tr,
		res.validateEntity,
		res.onAddPrepare,
		res.shouldSoftDelete,
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
}

func (s *teacherService) validateEntity(teacher entities.Teacher) error {
	if err := teacher.ValidateEmail(); err != nil {
		return err
	}

	return nil
}
func (s *teacherService) onAddPrepare(ctx context.Context, teacher *entities.Teacher) error {
	return nil
}
func (s *teacherService) shouldSoftDelete(teacher *entities.Teacher) bool {
	return true
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

	s.repository.Clear(ctx)
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
