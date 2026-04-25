package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/services/course/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type StudentCourseService interface {
	platform.BaseService[entities.StudentCourse]
}

func NewStudentCourseService(
	scr repositories.StudentCourseRepository,
	sr repositories.StudentRepository,
	cr repositories.CourseRepository,
) StudentCourseService {
	sc := platform.NewServiceConfig(
		"StudentCourseService",
		filepath.Join("data", "student_courses.json"),
		entities.StudentCourse{}.EntityName(),
	)

	res := &studentCourseService{
		repository:        scr,
		studentRepository: sr,
		courseRepository:  cr,
	}

	res.BaseService = platform.NewBaseService(sc, scr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.StudentCourse]{
			platform.OnAddPrepare: res.onAddPrepare,
		},
	)

	res.logger = res.GetLogger()
	return res
}

type studentCourseService struct {
	platform.BaseService[entities.StudentCourse]
	repository        repositories.StudentCourseRepository
	studentRepository repositories.StudentRepository
	courseRepository  repositories.CourseRepository
	logger            logger.Logger
}

func (s *studentCourseService) onAddPrepare(
	ctx context.Context, sc *entities.StudentCourse,
) error {
	sc.ID = uuid.New()
	return nil
}

func (s *studentCourseService) Seed(ctx context.Context) error {
	type seedItem struct {
		StudentName string `json:"student_name"`
		CourseName  string `json:"course_name"`
	}

	var mappings []seedItem
	if err := jsonutil.ReadFileTo(filepath.Join("data", "student_courses.json"), &mappings); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load student course seed data: %w", err), logger.ServiceValidationFailed,
		)
	}

	var lastError error

	for _, item := range mappings {
		student := s.studentRepository.FindFirstByName(ctx, item.StudentName)
		if student == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("student %q not found", item.StudentName), logger.ServiceValidationFailed,
			)
			continue
		}

		course := s.courseRepository.FindFirstByName(ctx, item.CourseName)
		if course == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("course %q not found", item.CourseName), logger.ServiceValidationFailed,
			)
			continue
		}

		scEntity := entities.StudentCourse{
			StudentID: student.ID,
			CourseID:  course.ID,
		}

		_, err := s.Add(ctx, scEntity)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to add %s: %w", scEntity, err), logger.ServiceValidationFailed,
			)
			continue
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d student course mappings processed from seed", len(mappings)), logger.ServiceOperationSuccess,
	)

	return lastError
}
