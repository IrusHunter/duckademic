package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/employees/entities"
	"github.com/IrusHunter/duckademic/services/employees/repositories"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

// AcademicRankService provides operations to initialize and manage academic ranks.
type AcademicRankService interface {
	// Seed clears existing academic ranks data and initializes it from a JSON file.
	Seed(context.Context) error
	// Add validates and inserts a new AcademicRank into the repository and returns it, or an error if it fails.
	Add(context.Context, entities.AcademicRank) (entities.AcademicRank, error)
	// GetAll returns a slice with all academic ranks.
	GetAll(context.Context) []entities.AcademicRank
	// Delete removes the AcademicRank with the specified ID from the repository.
	Delete(context.Context, uuid.UUID) error
	// Update updates the AcademicRank with the specified ID and returns the updated entity.
	Update(context.Context, uuid.UUID, entities.AcademicRank) (entities.AcademicRank, error)
	// FindByID returns a pointer to the academic rank from repository with the given id.
	FindByID(context.Context, uuid.UUID) *entities.AcademicRank
}

// NewAcademicRankService creates a new AcademicRankService instance.
//
// It requires a academic rank repository (arr).
func NewAcademicRankService(arr repositories.AcademicRankRepository) AcademicRankService {
	return &academicRankService{repository: arr}
}

type academicRankService struct {
	repository repositories.AcademicRankRepository
}

func (s *academicRankService) validateEntity(academicRank entities.AcademicRank) error {
	if err := academicRank.ValidateTitle(); err != nil {
		return err
	}

	return nil
}

func (s *academicRankService) Add(ctx context.Context, academicRank entities.AcademicRank,
) (entities.AcademicRank, error) {
	if err := s.validateEntity(academicRank); err != nil {
		return entities.AcademicRank{},
			fmt.Errorf("%s failed validation: %w", academicRank.String(), err)
	}

	slug := slug.Make(academicRank.Title)
	if other := s.repository.FindBySlug(ctx, slug); other != nil {
		return entities.AcademicRank{},
			fmt.Errorf("academic rank with slug %q already exists", slug)
	}
	academicRank.ID = uuid.New()
	academicRank.Slug = slug

	ar, err := s.repository.Add(ctx, academicRank)
	if err != nil {
		return entities.AcademicRank{},
			fmt.Errorf("failed to add %s to repository: %w", academicRank.String(), err)
	}
	return ar, nil
}
func (s *academicRankService) Seed(ctx context.Context) error {
	ranks := []entities.AcademicRank{}
	if err := jsonutil.ReadFileTo(filepath.Join("data", "academic_ranks.json"), &ranks); err != nil {
		return fmt.Errorf("failed to load academic ranks seed data: %w", err)
	}

	s.repository.Clear(ctx)
	var lastError error
	for _, academicRank := range ranks {
		_, err := s.Add(ctx, academicRank)
		if err != nil {
			lastError = fmt.Errorf("failed to add academicRank %q: %w", academicRank.Title, err)
		}
	}

	return lastError
}
func (s *academicRankService) GetAll(ctx context.Context) []entities.AcademicRank {
	return s.repository.GetAll(ctx)
}
func (s *academicRankService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete academic rank with id %q from repository: %w", id, err)
	}
	return nil
}
func (s *academicRankService) Update(ctx context.Context, id uuid.UUID, rank entities.AcademicRank,
) (entities.AcademicRank, error) {
	if err := s.validateEntity(rank); err != nil {
		return entities.AcademicRank{}, fmt.Errorf("%s failed validation: %w", rank.String(), err)
	}

	updatedR, err := s.repository.Update(ctx, id, rank)
	if err != nil {
		rank.ID = id
		return entities.AcademicRank{}, fmt.Errorf("failed to update %q in repository: %w", rank.String(), err)
	}

	return updatedR, nil
}
func (s *academicRankService) FindByID(ctx context.Context, id uuid.UUID) *entities.AcademicRank {
	return s.repository.FindByID(ctx, id)
}
