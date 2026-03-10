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
			fmt.Errorf("academic rank %q failed validation: %w", academicRank.Title, err)
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
			fmt.Errorf("failed to add academic rank %q to repository: %w", academicRank.Title, err)
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
