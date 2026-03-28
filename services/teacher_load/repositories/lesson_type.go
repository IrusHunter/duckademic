package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/teacher_load/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LessonTypeRepository interface {
	platform.BaseRepository[entities.LessonType]
	FindBySlug(ctx context.Context, slug string) *entities.LessonType
	FindFirstByName(ctx context.Context, name string) *entities.LessonType
	ExternalUpdate(context.Context, uuid.UUID, entities.LessonType) (entities.LessonType, error)
}

func NewLessonTypeRepository(db *sqlx.DB) LessonTypeRepository {
	config := platform.NewRepositoryConfig(
		"LessonTypeRepository",
		entities.LessonType{}.TableName(),
		"lesson_type",
		[]string{"id", "slug", "name"},
		[]string{""},
		[]string{"created_at", "updated_at"},
	)

	ltr := &lessonTypeRepository{
		BaseRepository: platform.NewBaseRepository[entities.LessonType](config, db),
		db:             db,
	}
	ltr.logger = ltr.GetLogger()

	return ltr
}

type lessonTypeRepository struct {
	platform.BaseRepository[entities.LessonType]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *lessonTypeRepository) FindBySlug(ctx context.Context, slug string) *entities.LessonType {
	return r.FindFirstBy(ctx, "slug", slug)
}
func (r *lessonTypeRepository) FindFirstByName(ctx context.Context, name string) *entities.LessonType {
	return r.FindFirstBy(ctx, "name", name)
}
func (r *lessonTypeRepository) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	lessonType entities.LessonType,
) (entities.LessonType, error) {
	return r.UpdateFields(ctx, id, []string{"slug", "name"}, lessonType)
}
