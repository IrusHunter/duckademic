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

type TeacherCourseService interface {
	platform.BaseService[entities.TeacherCourse]
}

func NewTeacherCourseService(
	tcr repositories.TeacherCourseRepository,
	tr repositories.TeacherRepository,
	cr repositories.CourseRepository,
) TeacherCourseService {
	sc := platform.NewServiceConfig(
		"TeacherCourseService",
		filepath.Join("data", "teacher_courses.json"),
		"teacher_course",
	)

	res := &teacherCourseService{
		repository:        tcr,
		teacherRepository: tr,
		courseRepository:  cr,
	}

	res.BaseService = platform.NewBaseService(sc, tcr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.TeacherCourse]{
			platform.OnAddPrepare: res.onAddPrepare,
		},
	)

	res.logger = res.GetLogger()
	return res
}

type teacherCourseService struct {
	platform.BaseService[entities.TeacherCourse]
	repository        repositories.TeacherCourseRepository
	teacherRepository repositories.TeacherRepository
	courseRepository  repositories.CourseRepository
	logger            logger.Logger
}

func (s *teacherCourseService) onAddPrepare(
	ctx context.Context,
	tc *entities.TeacherCourse,
) error {
	tc.ID = uuid.New()
	return nil
}

func (s *teacherCourseService) Seed(ctx context.Context) error {
	type seedItem struct {
		TeacherName string `json:"teacher_name"`
		CourseName  string `json:"course_name"`
	}

	var items []seedItem
	if err := jsonutil.ReadFileTo(filepath.Join("data", "teacher_courses.json"), &items); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load teacher course seed data: %w", err),
			logger.ServiceValidationFailed,
		)
	}

	var lastError error

	for _, item := range items {
		teacher := s.teacherRepository.FindByName(ctx, item.TeacherName)
		if teacher == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("teacher %q not found", item.TeacherName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		course := s.courseRepository.FindFirstByName(ctx, item.CourseName)
		if course == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("course %q not found", item.CourseName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		tc := entities.TeacherCourse{
			TeacherID: teacher.ID,
			CourseID:  course.ID,
		}

		_, err := s.Add(ctx, tc)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to add %s: %w", tc, err),
				logger.ServiceValidationFailed,
			)
			continue
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d teacher course mappings processed from seed", len(items)),
		logger.ServiceOperationSuccess,
	)

	return lastError
}
