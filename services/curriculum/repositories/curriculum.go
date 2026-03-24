package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/curriculum/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type CurriculumRepository interface {
	platform.BaseRepository[entities.Curriculum]
	FindBySlug(context.Context, string) *entities.Curriculum
	FindFirstByName(ctx context.Context, name string) *entities.Curriculum
}

func NewCurriculumRepository(db *sqlx.DB) CurriculumRepository {
	config := platform.NewRepositoryConfig(
		"CurriculumRepository",
		entities.Curriculum{}.TableName(),
		"curriculum",
		[]string{"id", "slug", "name", "duration_years", "effective_from", "effective_to"},
		[]string{"name", "duration_years", "effective_from", "effective_to"},
		[]string{"created_at", "updated_at"},
	)

	cr := &curriculumRepository{
		BaseRepository: platform.NewBaseRepository[entities.Curriculum](config, db),
		db:             db,
	}
	cr.logger = cr.GetLogger()

	return cr
}

type curriculumRepository struct {
	platform.BaseRepository[entities.Curriculum]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *curriculumRepository) FindBySlug(ctx context.Context, slug string) *entities.Curriculum {
	return r.FindFirstBy(ctx, "slug", slug)
}
func (r *curriculumRepository) FindFirstByName(ctx context.Context, name string) *entities.Curriculum {
	return r.FindFirstBy(ctx, "name", name)
}
