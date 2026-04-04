package services

import (
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/entities"
)

type ValidationService interface {
	ValidateTeachers([]entities.Teacher) error
	ValidateDisciplines([]entities.Discipline) error
}

func NewValidationService() ValidationService {
	return &validationService{}
}

type validationService struct{}

func (s *validationService) ValidateTeachers(teachers []entities.Teacher) error {
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

func (s *validationService) ValidateDisciplines(disciplines []entities.Discipline) error {
	if len(disciplines) == 0 {
		return fmt.Errorf("disciplines list cannot be empty")
	}

	seenIDs := make(map[string]struct{})

	for i, d := range disciplines {
		if err := d.ValidateName(); err != nil {
			return fmt.Errorf("failed name validation for discipline at index %d: %w", i, err)
		}

		if _, exists := seenIDs[d.ID.String()]; exists {
			return fmt.Errorf("duplicate discipline ID found: %s", d.ID)
		}
		seenIDs[d.ID.String()] = struct{}{}
	}

	return nil
}
