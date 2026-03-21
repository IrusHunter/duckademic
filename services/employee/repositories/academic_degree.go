package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/employees/entities"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

// AcademicDegreeRepository represents a storage for academic rank entities.
type AcademicDegreeRepository interface {
	platform.BaseRepository[entities.AcademicDegree]
	// FindBySlug returns a pointer to the academic degree from database with the given slug.
	FindBySlug(context.Context, string) *entities.AcademicDegree
	FindByTitle(context.Context, string) *entities.AcademicDegree
}

// NewAcademicRankRepository creates a new AcademicDegreeRepository instance.
//
// It requires a database connection (db).
func NewAcademicDegreeRepository(db *sqlx.DB) AcademicDegreeRepository {
	config := platform.NewRepositoryConfig("AcademicDegreeRepository", "academic_degrees", "academic degree",
		[]string{"id", "slug", "title"}, []string{"title"}, []string{"created_at", "updated_at"},
	)
	return &academicDegreeRepository{
		BaseRepository: platform.NewBaseRepository[entities.AcademicDegree](config, db),
		db:             db,
	}
}

type academicDegreeRepository struct {
	platform.BaseRepository[entities.AcademicDegree]
	db *sqlx.DB
}

func (r *academicDegreeRepository) FindBySlug(ctx context.Context, slug string) *entities.AcademicDegree {
	return r.FindFirstBy(ctx, "slug", slug)
}
func (r *academicDegreeRepository) FindByTitle(ctx context.Context, title string) *entities.AcademicDegree {
	return r.FindFirstBy(ctx, "title", title)
}
