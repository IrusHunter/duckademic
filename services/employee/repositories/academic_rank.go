package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/employees/entities"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

// AcademicRankRepository represents a storage for academic rank entities.
type AcademicRankRepository interface {
	platform.BaseRepository[entities.AcademicRank]
	// FindBySlug returns a pointer to the academic rank from database with the given slug.
	FindBySlug(context.Context, string) *entities.AcademicRank
}

// NewAcademicRankRepository creates a new AcademicRankRepository instance.
//
// It requires a database connection (db).
func NewAcademicRankRepository(db *sqlx.DB) AcademicRankRepository {
	config := platform.NewRepositoryConfig("AcademicRankRepository", "academic_ranks", "academic rank",
		[]string{"id", "slug", "title"}, []string{"id", "slug", "title", "created_at", "updated_at"}, []string{"title"},
		[]string{"created_at", "updated_at"},
	)
	return &academicRankRepository{
		BaseRepository: platform.NewBaseRepository[entities.AcademicRank](config, db),
		db:             db,
	}
}

type academicRankRepository struct {
	platform.BaseRepository[entities.AcademicRank]
	db *sqlx.DB
}

func (r *academicRankRepository) FindBySlug(ctx context.Context, slug string) *entities.AcademicRank {
	return r.FindFirstBy(ctx, "slug", slug)
}
