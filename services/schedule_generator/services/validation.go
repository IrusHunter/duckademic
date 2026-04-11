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
	ValidateGroupCohorts([]entities.GroupCohort) error
	ValidateClassrooms([]entities.Classroom) error
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
func (s *validationService) ValidateGroupCohorts(cohorts []entities.GroupCohort) error {
	if len(cohorts) == 0 {
		return fmt.Errorf("group cohort list cannot be empty")
	}

	seenCohortIDs := make(map[string]struct{})

	for i, cohort := range cohorts {
		if err := cohort.ValidateGroups(); err != nil {
			return fmt.Errorf("failed groups validation for cohort at index %d: %w", i, err)
		}

		if _, exists := seenCohortIDs[cohort.ID.String()]; exists {
			return fmt.Errorf("duplicate group cohort ID found: %s", cohort.ID)
		}
		seenCohortIDs[cohort.ID.String()] = struct{}{}

		seenGroupIDs := make(map[string]struct{})

		for j, g := range cohort.Groups {
			if err := g.ValidateName(); err != nil {
				return fmt.Errorf("failed name validation for group at index %d in cohort %d: %w", j, i, err)
			}

			if err := g.ValidateConnectedGroups(); err != nil {
				return fmt.Errorf("failed connected groups validation for group at index %d in cohort %d: %w", j, i, err)
			}

			if err := g.ValidateStudentCount(); err != nil {
				return fmt.Errorf("failed student count validation for group at index %d in cohort %d: %w", j, i, err)
			}

			if _, exists := seenGroupIDs[g.ID.String()]; exists {
				return fmt.Errorf("duplicate student group ID found in cohort %d: %s", i, g.ID)
			}
			seenGroupIDs[g.ID.String()] = struct{}{}
		}
	}

	return nil
}
func (s *validationService) ValidateClassrooms(classrooms []entities.Classroom) error {
	if len(classrooms) == 0 {
		return fmt.Errorf("classrooms list cannot be empty")
	}

	seenIDs := make(map[string]struct{})

	for i, c := range classrooms {
		if err := c.ValidateNumber(); err != nil {
			return fmt.Errorf("failed number validation for classroom at index %d: %w", i, err)
		}

		if err := c.ValidateCapacity(); err != nil {
			return fmt.Errorf("failed capacity validation for classroom at index %d: %w", i, err)
		}

		if _, exists := seenIDs[c.ID.String()]; exists {
			return fmt.Errorf("duplicate classroom ID found: %s", c.ID)
		}
		seenIDs[c.ID.String()] = struct{}{}
	}

	return nil
}
