package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type StudentService interface {
	platform.BaseService[entities.Student]
}

func NewStudentService(tr repositories.StudentRepository, eb events.EventBus) StudentService {
	sc := platform.NewServiceConfig("StudentService", filepath.Join("data", "students.json"), "student")

	res := &studentService{
		repository: tr,
	}
	res.BaseService = platform.NewBaseService(sc, tr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Student]{
			platform.HardDeleteCheck: res.hardDeleteCheck,
		},
	)
	res.logger = res.GetLogger()

	eb.Subscribe(contextutil.SetTraceID(context.Background()), string(events.StudentRT), res.eventHandler)

	return res
}

type studentService struct {
	platform.BaseService[entities.Student]
	repository repositories.StudentRepository
	logger     logger.Logger
}

func (s *studentService) hardDeleteCheck(ctx context.Context, student *entities.Student) error {
	return fmt.Errorf("plug")
}
func (s *studentService) eventHandler(ctx context.Context, b []byte) {
	student, err := events.FromByteConvertor[events.StudentRE](b)
	if err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "StudentRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "StudentRTHandler",
		fmt.Sprintf("received %s", student), logger.EventDataReceived,
	)

	trueS := entities.Student{
		ID:   student.ID,
		Slug: student.Slug,
		Name: student.Name,
	}

	switch student.Event {
	case events.EntityCreated:
		s.Add(ctx, trueS)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, student.ID, trueS)
	case events.EntityDeleted:
		s.Delete(ctx, student.ID)
	}
}

func (s *studentService) Seed(ctx context.Context) error {
	students := []entities.Student{}
	if err := jsonutil.ReadFileTo(filepath.Join("data", "students.json"), &students); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load students seed data: %w", err), logger.ServiceDataFetchFailed,
		)
	}

	var lastError error
	for _, student := range students {
		trueS := s.repository.FindFirstByName(ctx, student.Name)
		if trueS == nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("student with name %q not found", student.Name), logger.ServiceDataFetchFailed,
			)
			continue
		}

		_, err := s.Update(ctx, trueS.ID, student)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to update %s: %w", student, err), logger.ServiceRepositoryFailed,
			)
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d students updated successfully", len(students)), logger.ServiceOperationSuccess,
	)
	return lastError
}

func (s *studentService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	student entities.Student,
) (entities.Student, error) {
	updatedS, err := s.repository.ExternalUpdate(ctx, id, student)
	if err != nil {
		return entities.Student{}, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "ExternalUpdate",
			err, logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedS), logger.ServiceOperationSuccess)
	return updatedS, nil
}
