package services

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/services/course/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type CourseService interface {
	platform.BaseService[entities.Course]
	GetStudentCoursePage(ctx context.Context, userID uuid.UUID) ([]StudentCoursePage, error)
}

func NewCourseService(
	cr repositories.CourseRepository,
	tr repositories.TeacherRepository,
	taskStudentRepository repositories.TaskStudentRepository,
	studentCourseRepository repositories.StudentCourseRepository,
	eb events.EventBus,
) CourseService {
	sc := platform.NewServiceConfig("CourseService", filepath.Join("data", "courses.json"), "course")

	res := &courseService{
		repository:              cr,
		teacherRepository:       tr,
		taskStudentRepository:   taskStudentRepository,
		studentCourseRepository: studentCourseRepository,
	}
	res.BaseService = platform.NewBaseService(sc, cr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Course]{},
	)
	res.logger = res.GetLogger()

	eb.Subscribe(contextutil.SetTraceID(context.Background()), string(events.DisciplineRT), res.eventHandler)

	return res
}

type courseService struct {
	platform.BaseService[entities.Course]
	repository              repositories.CourseRepository
	teacherRepository       repositories.TeacherRepository
	taskStudentRepository   repositories.TaskStudentRepository
	studentCourseRepository repositories.StudentCourseRepository
	logger                  logger.Logger
}

func (s *courseService) eventHandler(ctx context.Context, b []byte) {
	cr, err := events.FromByteConvertor[events.DisciplineRE](b)
	if err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "CourseRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "CourseRTHandler",
		fmt.Sprintf("received %s", cr), logger.EventDataReceived,
	)

	trueCR := entities.Course{
		ID:   cr.ID,
		Slug: cr.Slug,
		Name: cr.Name,
	}

	switch cr.Event {
	case events.EntityCreated:
		s.Add(ctx, trueCR)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, cr.ID, trueCR)
	case events.EntityDeleted:
		s.Delete(ctx, cr.ID)
	}
}

func (s *courseService) Seed(ctx context.Context) error {
	type seedCourse struct {
		Name        string `json:"name"`
		ManagerName string `json:"manager_name"`
		Description string `json:"description"`
	}

	courses := []seedCourse{}
	if err := jsonutil.ReadFileTo(filepath.Join("data", "courses.json"), &courses); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load courses seed data: %w", err), logger.ServiceDataFetchFailed,
		)
	}

	var lastError error
	for _, course := range courses {
		existing := s.repository.FindFirstByName(ctx, course.Name)
		if existing == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("course with name %q not found", course.Name), logger.ServiceDataFetchFailed,
			)
			continue
		}

		manager := s.teacherRepository.FindByName(ctx, course.ManagerName)
		if manager == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("manager with name %q not found", course.ManagerName), logger.ServiceDataFetchFailed,
			)
			continue
		}

		updated := entities.Course{
			Name:        course.Name,
			Description: &course.Description,
			ManagerID:   &manager.ID,
		}

		_, err := s.Update(ctx, existing.ID, updated)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to update %s: %w", updated, err), logger.ServiceRepositoryFailed,
			)
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d courses updated successfully", len(courses)), logger.ServiceOperationSuccess,
	)
	return lastError
}

func (s *courseService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	course entities.Course,
) (entities.Course, error) {
	updatedCR, err := s.repository.ExternalUpdate(ctx, id, course)
	if err != nil {
		return entities.Course{}, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "ExternalUpdate",
			err, logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedCR), logger.ServiceOperationSuccess)
	return updatedCR, nil
}

type StudentCoursePage struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`

	TeacherName string `json:"teacher_name"`

	AverageMark      float64   `json:"average_mark"`
	AssignmentsCount int       `json:"assignments_count"`
	StudentCount     int       `json:"student_count"`
	UpcomingDeadline time.Time `json:"upcoming_deadline"`
}

func (s *courseService) GetStudentCoursePage(ctx context.Context, userID uuid.UUID) ([]StudentCoursePage, error) {
	courses, err := s.repository.GetFullCourses(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.formStudentCoursePages(ctx, userID, courses)
}

func (s *courseService) formStudentCoursePages(
	ctx context.Context, studentID uuid.UUID, courses []entities.Course,
) ([]StudentCoursePage, error) {
	result := make([]StudentCoursePage, 0, len(courses))
	now := time.Now()

	for _, course := range courses {
		taskStudents, err := s.taskStudentRepository.GetTasksForStudentInCourse(ctx, studentID, course.ID)
		if err != nil {
			return nil, err
		}

		page := StudentCoursePage{
			ID:          course.ID,
			Name:        course.Name,
			Description: valueOrEmpty(course.Description),
		}

		if course.Manager != nil {
			page.TeacherName = course.Manager.Name
		}

		page.AssignmentsCount = len(taskStudents)

		var (
			marksSum    float64
			marksCount  int
			nearestTime *time.Time
		)

		for _, ts := range taskStudents {
			if ts.Mark != nil {
				marksSum += *ts.Mark
				marksCount++
			}

			if ts.Task != nil &&
				ts.Task.Deadline.After(now) {

				if nearestTime == nil ||
					ts.Task.Deadline.Before(*nearestTime) {

					deadline := ts.Task.Deadline
					nearestTime = &deadline
				}
			}
		}

		if marksCount > 0 {
			page.AverageMark = marksSum / float64(marksCount)
		}

		if nearestTime != nil {
			page.UpcomingDeadline = *nearestTime
		}

		students, err := s.studentCourseRepository.GetStudentsForCourse(ctx, course.ID)
		if err != nil {
			return nil, err
		}

		page.StudentCount = len(students)

		result = append(result, page)
	}

	return result, nil
}

func valueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
