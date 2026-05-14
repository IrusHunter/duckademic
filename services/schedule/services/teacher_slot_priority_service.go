package services

import (
	"context"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/repositories"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

// TeacherSlotPriorityService provides operations for teacher slot priority management.
type TeacherSlotPriorityService interface {
	platform.BaseService[entities.TeacherSlotPriority]
}

// NewTeacherSlotPriorityService creates a new TeacherSlotPriorityService instance.
func NewTeacherSlotPriorityService(
	r repositories.TeacherSlotPriorityRepository,
) TeacherSlotPriorityService {
	sc := platform.NewServiceConfig(
		"TeacherSlotPriorityService",
		filepath.Join("data", "teacher_slot_priorities.json"),
		entities.TeacherSlotPriority{}.EntityName(),
	)

	s := &teacherSlotPriorityService{
		repository: r,
	}

	s.BaseService = platform.NewBaseService(sc, r,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.TeacherSlotPriority]{
			platform.OnAddPrepare:   s.onAddPrepare,
			platform.ValidateEntity: s.validateEntity,
		},
	)

	return s
}

type teacherSlotPriorityService struct {
	platform.BaseService[entities.TeacherSlotPriority]
	repository repositories.TeacherSlotPriorityRepository
}

func (s *teacherSlotPriorityService) validateEntity(ctx context.Context, tsp *entities.TeacherSlotPriority) error {
	if err := tsp.ValidatePriority(); err != nil {
		return err
	}

	return nil
}
func (s *teacherSlotPriorityService) onAddPrepare(ctx context.Context, tsp *entities.TeacherSlotPriority) error {
	tsp.ID = uuid.New()

	return nil
}
