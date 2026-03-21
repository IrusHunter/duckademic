package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// AcademicRankRepository represents a storage for academic rank entities.
type AcademicRankRepository interface {
	platform.BaseRepository[entities.AcademicRank]
	// FindBySlug returns a pointer to the academic rank from database with the given slug.
	FindBySlug(context.Context, string) *entities.AcademicRank
	FindByTitle(context.Context, string) *entities.AcademicRank
	ExternalUpdate(context.Context, uuid.UUID, entities.AcademicRank) (entities.AcademicRank, error)
}

// NewAcademicRankRepository creates a new AcademicRankRepository instance.
//
// It requires a database connection (db).
func NewAcademicRankRepository(db *sqlx.DB) AcademicRankRepository {
	config := platform.NewRepositoryConfig("AcademicRankRepository", entities.AcademicRank{}.TableName(),
		"academic rank", []string{"id", "slug", "title"}, []string{"priority"}, []string{"created_at", "updated_at"},
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
func (r *academicRankRepository) FindByTitle(ctx context.Context, title string) *entities.AcademicRank {
	return r.FindFirstBy(ctx, "title", title)
}
func (r *academicRankRepository) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	academicRank entities.AcademicRank,
) (entities.AcademicRank, error) {
	return r.UpdateFields(ctx, id, []string{"title"}, academicRank)
}
