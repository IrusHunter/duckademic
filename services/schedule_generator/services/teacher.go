package services

import (
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
)

type TeacherService interface {
	ValidateTeachers([]entities.Teacher) error
}

func NewTeacherService() TeacherService {
	return &teacherService{
		logger: logger.NewLogger("TeacherService.txt", "TeacherService"),
	}
}

type teacherService struct {
	logger logger.Logger
}

func (s *teacherService) ValidateTeachers(teachers []entities.Teacher) error {
	if len(teachers) == 0 {
		return fmt.Errorf("teachers list cannot be empty")
	}

	seenIDs := make(map[string]struct{})

	for i, t := range teachers {
		if err := t.ValidateName(); err != nil {
			return fmt.Errorf("failed name validation for teacher at index %d: %w", i, err)
		}

		if _, exists := seenIDs[t.ID.String()]; exists {
			return fmt.Errorf("duplicate teacher ID found: %s", t.ID)
		}
		seenIDs[t.ID.String()] = struct{}{}
	}

	return nil
}
