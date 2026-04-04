package services

import (
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/entities"
)

type ValidationService interface {
	ValidateTeachers([]entities.Teacher) error
	ValidateDisciplines([]entities.Discipline) error
	ValidateLessonTypeRequests([]entities.LessonTypeRequest) error
	ValidateLessonTypeAssignments([]entities.LessonTypeAssignment) error
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
func (s *validationService) ValidateLessonTypeRequests(lessonTypes []entities.LessonTypeRequest) error {
	if len(lessonTypes) == 0 {
		return fmt.Errorf("lesson type list cannot be empty")
	}

	seenIDs := make(map[string]struct{})

	for i, lt := range lessonTypes {
		if err := lt.ValidateName(); err != nil {
			return fmt.Errorf("failed name validation for lesson type at index %d: %w", i, err)
		}

		if _, exists := seenIDs[lt.ID.String()]; exists {
			return fmt.Errorf("duplicate lesson type ID found: %s", lt.ID)
		}
		seenIDs[lt.ID.String()] = struct{}{}
	}

	return nil
}
func (s *validationService) ValidateLessonTypeAssignments(assignments []entities.LessonTypeAssignment) error {
	if len(assignments) == 0 {
		return fmt.Errorf("lesson type assignment list cannot be empty")
	}

	seenIDs := make(map[string]struct{})

	for i, a := range assignments {
		if err := a.ValidateRequiredHours(); err != nil {
			return fmt.Errorf("failed required hours validation for assignment at index %d: %w", i, err)
		}

		if _, exists := seenIDs[a.ID.String()]; exists {
			return fmt.Errorf("duplicate lesson type assignment ID found: %s", a.ID)
		}
		seenIDs[a.ID.String()] = struct{}{}
	}

	return nil
}
