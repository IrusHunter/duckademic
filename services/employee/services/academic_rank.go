package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/employees/entities"
	"github.com/IrusHunter/duckademic/services/employees/repositories"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

// AcademicRankService provides operations to initialize and manage academic ranks.
type AcademicRankService interface {
	platform.BaseService[entities.AcademicRank]
}

// NewAcademicRankService creates a new AcademicRankService instance.
//
// It requires a academic rank repository (arr).
func NewAcademicRankService(arr repositories.AcademicRankRepository) AcademicRankService {
	sc := platform.NewServiceConfig("AcademicRankService", filepath.Join("data", "academic_ranks.json"), "academic_rank")

	res := &academicRankService{
		repository: arr,
	}
	res.BaseService = platform.NewBaseService(sc, arr, res.validateEntity, res.onAddPrepare,
		func(ar *entities.AcademicRank) bool { return false },
	)

	return res
}

type academicRankService struct {
	platform.BaseService[entities.AcademicRank]
	repository repositories.AcademicRankRepository
}

func (s *academicRankService) validateEntity(academicRank entities.AcademicRank) error {
	if err := academicRank.ValidateTitle(); err != nil {
		return err
	}

	return nil
}

func (s *academicRankService) onAddPrepare(ctx context.Context, academicRank *entities.AcademicRank) error {
	slug := slug.Make(academicRank.Title)
	if other := s.repository.FindBySlug(ctx, slug); other != nil {
		return fmt.Errorf("academic rank with slug %q already exists", slug)
	}
	academicRank.ID = uuid.New()
	academicRank.Slug = slug

	return nil
}
