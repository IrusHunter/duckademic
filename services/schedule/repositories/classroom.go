package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ClassroomRepository interface {
	platform.BaseRepository[entities.Classroom]
	FindBySlug(ctx context.Context, slug string) *entities.Classroom
	FindFirstByNumber(ctx context.Context, number string) *entities.Classroom
	ExternalUpdate(ctx context.Context, id uuid.UUID, classroom entities.Classroom) (entities.Classroom, error)
}

func NewClassroomRepository(db *sqlx.DB) ClassroomRepository {
	config := platform.NewRepositoryConfig(
		"ClassroomRepository",
		entities.Classroom{}.TableName(),
		"classroom",
		[]string{"id", "slug", "number", "capacity"},
		[]string{},
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

func (r *classroomRepository) FindFirstByNumber(ctx context.Context, number string) *entities.Classroom {
	return r.FindFirstBy(ctx, "number", number)
}

func (r *classroomRepository) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	classroom entities.Classroom,
) (entities.Classroom, error) {
	return r.UpdateFields(ctx, id, []string{"number", "capacity"}, classroom)
}
