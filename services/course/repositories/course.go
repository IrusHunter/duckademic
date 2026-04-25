package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CourseRepository interface {
	platform.BaseRepository[entities.Course]
	FindFirstByName(ctx context.Context, name string) *entities.Course
	ExternalUpdate(context.Context, uuid.UUID, entities.Course) (entities.Course, error)
}

func NewCourseRepository(db *sqlx.DB) CourseRepository {
	config := platform.NewRepositoryConfig(
		"CourseRepository",
		entities.Course{}.TableName(),
		entities.Course{}.EntityName(),
		[]string{"id", "slug", "name"},
		[]string{"manager_id", "description"},
		[]string{"created_at", "updated_at"},
	)

	cr := &courseRepository{
		BaseRepository: platform.NewBaseRepository[entities.Course](config, db),
		db:             db,
	}
	cr.logger = cr.GetLogger()

	return cr
}

type courseRepository struct {
	platform.BaseRepository[entities.Course]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *courseRepository) FindFirstByName(ctx context.Context, name string) *entities.Course {
	return r.FindFirstBy(ctx, "name", name)
}

func (r *courseRepository) ExternalUpdate(
	ctx context.Context, id uuid.UUID, course entities.Course,
) (entities.Course, error) {
	return r.UpdateFields(ctx, id, []string{"slug", "name"}, course)
}
