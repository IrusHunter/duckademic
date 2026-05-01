package services

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

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
	GetCourseProgress(context.Context, uuid.UUID) ([]CourseProgress, error)
}

func NewStudentCourseService(
	scr repositories.StudentCourseRepository,
	sr repositories.StudentRepository,
	cr repositories.CourseRepository,
	tsr repositories.TaskStudentRepository,
) StudentCourseService {
	sc := platform.NewServiceConfig(
		"StudentCourseService",
		filepath.Join("data", "student_courses.json"),
		entities.StudentCourse{}.EntityName(),
	)

	res := &studentCourseService{
		repository:            scr,
		studentRepository:     sr,
		courseRepository:      cr,
		taskStudentRepository: tsr,
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
	repository            repositories.StudentCourseRepository
	studentRepository     repositories.StudentRepository
	courseRepository      repositories.CourseRepository
	taskStudentRepository repositories.TaskStudentRepository
	logger                logger.Logger
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

type CourseProgress struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	CompleteRate     float32   `json:"complete_rate"`
	CompleteAccuracy float32   `json:"complete_accuracy"`
}

func (s *studentCourseService) GetCourseProgress(ctx context.Context, studentID uuid.UUID) ([]CourseProgress, error) {
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup
	var mu sync.Mutex

	var result []CourseProgress
	var lastError error

	courses, err := s.repository.GetCoursesForStudent(ctx, studentID)
	if err != nil {
		return nil, s.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetCourseProgress",
			err,
			logger.ServiceRepositoryFailed,
		)
	}

	for _, course := range courses {
		wg.Add(1)
		sem <- struct{}{}

		go func(course entities.Course) {
			defer wg.Done()
			defer func() { <-sem }()

			taskStudents, err := s.taskStudentRepository.GetTasksForStudentInCourse(
				ctx,
				studentID,
				course.ID,
			)
			if err != nil {
				mu.Lock()
				lastError = s.GetLogger().LogAndReturnError(
					contextutil.GetTraceID(ctx),
					"GetCourseProgress",
					err,
					logger.ServiceRepositoryFailed,
				)
				mu.Unlock()
				return
			}

			totalTasks := len(taskStudents)
			if totalTasks == 0 {
				return
			}

			var completed int
			var accuracy float32
			var graded int

			taskMap := make(map[uuid.UUID]entities.TaskStudent)
			for _, ts := range taskStudents {
				taskMap[ts.TaskID] = ts

				if ts.SubmissionTime != nil {
					completed++
				}

				if ts.Mark != nil {
					accuracy = (accuracy*float32(graded) + float32(*ts.Mark)/float32(ts.Task.MaxMark)) / float32(graded+1)
					graded++
				}
			}

			progress := CourseProgress{
				ID:               course.ID,
				Name:             course.Name,
				CompleteRate:     float32(completed) / float32(totalTasks),
				CompleteAccuracy: accuracy,
			}

			mu.Lock()
			result = append(result, progress)
			mu.Unlock()

		}(course)
	}

	wg.Wait()

	if lastError != nil {
		return result, lastError
	}

	return result, nil
}
