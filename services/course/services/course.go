package services

import (
	"context"
	"fmt"
	"path/filepath"

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
}

func NewCourseService(
	cr repositories.CourseRepository,
	tr repositories.TeacherRepository,
	eb events.EventBus,
) CourseService {
	sc := platform.NewServiceConfig("CourseService", filepath.Join("data", "courses.json"), "course")

	res := &courseService{
		repository:        cr,
		teacherRepository: tr,
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
	repository        repositories.CourseRepository
	teacherRepository repositories.TeacherRepository
	logger            logger.Logger
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
			Description: course.Description,
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
