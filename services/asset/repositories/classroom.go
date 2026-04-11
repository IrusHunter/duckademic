package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/asset/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type ClassroomRepository interface {
	platform.BaseRepository[entities.Classroom]
	FindBySlug(context.Context, string) *entities.Classroom
}

func NewClassroomRepository(db *sqlx.DB) ClassroomRepository {
	config := platform.NewRepositoryConfig(
		"ClassroomRepository",
		entities.Classroom{}.TableName(),
		"classroom",
		[]string{"id", "slug", "number", "capacity"},
		[]string{"number", "capacity"},
		[]string{"created_at", "updated_at"},
	)

	cr := &classroomRepository{
		BaseRepository: platform.NewBaseRepository[entities.Classroom](config, db),
		db:             db,
	}
	cr.logger = cr.GetLogger()

	return cr
}

type classroomRepository struct {
	platform.BaseRepository[entities.Classroom]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *classroomRepository) FindBySlug(ctx context.Context, slug string) *entities.Classroom {
	return r.FindFirstBy(ctx, "slug", slug)
}
