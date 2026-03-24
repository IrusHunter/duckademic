package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/curriculum/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type DisciplineRepository interface {
	platform.BaseRepository[entities.Discipline]
	FindBySlug(ctx context.Context, slug string) *entities.Discipline
	FindFirstByName(ctx context.Context, name string) *entities.Discipline
}

func NewDisciplineRepository(db *sqlx.DB) DisciplineRepository {
	config := platform.NewRepositoryConfig(
		"DisciplineRepository",
		entities.Discipline{}.TableName(),
		"discipline",
		[]string{"id", "slug", "name"},
		[]string{"name"},
		[]string{"created_at", "updated_at"},
	)

	dr := &disciplineRepository{
		BaseRepository: platform.NewBaseRepository[entities.Discipline](config, db),
		db:             db,
	}
	dr.logger = dr.GetLogger()

	return dr
}

type disciplineRepository struct {
	platform.BaseRepository[entities.Discipline]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *disciplineRepository) FindBySlug(ctx context.Context, slug string) *entities.Discipline {
	return r.FindFirstBy(ctx, "slug", slug)
}
func (r *disciplineRepository) FindFirstByName(ctx context.Context, name string) *entities.Discipline {
	return r.FindFirstBy(ctx, "name", name)
}
