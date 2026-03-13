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

// AcademicDegreeService provides operations to initialize and manage academic degrees.
type AcademicDegreeService interface {
	platform.BaseService[entities.AcademicDegree]
}

// NewAcademicDegreeService creates a new AcademicDegreeService instance.
//
// It requires a academic degree repository (adr).
func NewAcademicDegreeService(adr repositories.AcademicDegreeRepository) AcademicDegreeService {
	sc := platform.NewServiceConfig("AcademicDegreeService", filepath.Join("data", "academic_degrees.json"), "academic_degree")

	res := &academicDegreeService{
		repository: adr,
	}
	res.BaseService = platform.NewBaseService(sc, adr, res.validateEntity, res.onAddPrepare)

	return res
}

type academicDegreeService struct {
	platform.BaseService[entities.AcademicDegree]
	repository repositories.AcademicDegreeRepository
}

func (s *academicDegreeService) validateEntity(AcademicDegree entities.AcademicDegree) error {
	if err := AcademicDegree.ValidateTitle(); err != nil {
		return err
	}

	return nil
}

func (s *academicDegreeService) onAddPrepare(ctx context.Context, academicDegree *entities.AcademicDegree) error {
	slug := slug.Make(academicDegree.Title)
	if other := s.repository.FindBySlug(ctx, slug); other != nil {
		return fmt.Errorf("academic degree with slug %q already exists", slug)
	}
	academicDegree.ID = uuid.New()
	academicDegree.Slug = slug

	return nil
}
