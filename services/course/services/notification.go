package services

import (
	"context"
	"fmt"

	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/services/course/repositories"
	"github.com/google/uuid"
)

type NotificationService interface {
	GetForUser(ctx context.Context, userID uuid.UUID) ([]entities.Notification, error)
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
	NotifyNewTask(ctx context.Context, task entities.Task)
	NotifySubmission(ctx context.Context, taskID, studentID uuid.UUID)
	NotifyGrade(ctx context.Context, taskID, studentID uuid.UUID)
}

func NewNotificationService(
	notifRepo repositories.NotificationRepository,
	studentCourseRepo repositories.StudentCourseRepository,
	taskRepo repositories.TaskRepository,
	courseRepo repositories.CourseRepository,
	studentRepo repositories.StudentRepository,
) NotificationService {
	return &notificationService{
		notifRepo:         notifRepo,
		studentCourseRepo: studentCourseRepo,
		taskRepo:          taskRepo,
		courseRepo:        courseRepo,
		studentRepo:       studentRepo,
	}
}

type notificationService struct {
	notifRepo         repositories.NotificationRepository
	studentCourseRepo repositories.StudentCourseRepository
	taskRepo          repositories.TaskRepository
	courseRepo        repositories.CourseRepository
	studentRepo       repositories.StudentRepository
}

func (s *notificationService) GetForUser(ctx context.Context, userID uuid.UUID) ([]entities.Notification, error) {
	return s.notifRepo.GetForRecipient(ctx, userID)
}

func (s *notificationService) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	return s.notifRepo.MarkAllReadForRecipient(ctx, userID)
}

func (s *notificationService) NotifyNewTask(ctx context.Context, task entities.Task) {
	if task.PostType == "announcement" {
		return
	}
	students, err := s.studentCourseRepo.GetStudentsForCourse(ctx, task.CourseID)
	if err != nil || len(students) == 0 {
		return
	}
	for _, student := range students {
		n := entities.Notification{
			ID:          uuid.New(),
			RecipientID: student.ID,
			Type:        "new_task",
			Title:       fmt.Sprintf("New assignment: %s", task.Title),
			TaskID:      task.ID,
			CourseID:    task.CourseID,
		}
		_ = s.notifRepo.Create(ctx, n)
	}
}

func (s *notificationService) NotifySubmission(ctx context.Context, taskID, studentID uuid.UUID) {
	task := s.taskRepo.FindByID(ctx, taskID)
	if task == nil {
		return
	}
	course := s.courseRepo.FindByID(ctx, task.CourseID)
	if course == nil || course.ManagerID == nil {
		return
	}
	studentName := "A student"
	if student := s.studentRepo.FindByID(ctx, studentID); student != nil {
		studentName = student.Name
	}
	n := entities.Notification{
		ID:          uuid.New(),
		RecipientID: *course.ManagerID,
		Type:        "submission",
		Title:       fmt.Sprintf("%s submitted: %s", studentName, task.Title),
		TaskID:      task.ID,
		CourseID:    task.CourseID,
	}
	_ = s.notifRepo.Create(ctx, n)
}

func (s *notificationService) NotifyGrade(ctx context.Context, taskID, studentID uuid.UUID) {
	task := s.taskRepo.FindByID(ctx, taskID)
	if task == nil {
		return
	}
	n := entities.Notification{
		ID:          uuid.New(),
		RecipientID: studentID,
		Type:        "grade",
		Title:       fmt.Sprintf("Grade received for: %s", task.Title),
		TaskID:      task.ID,
		CourseID:    task.CourseID,
	}
	_ = s.notifRepo.Create(ctx, n)
}
