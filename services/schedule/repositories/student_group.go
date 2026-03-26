package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type StudentGroupRepository interface {
	platform.BaseRepository[entities.StudentGroup]
	FindBySlug(context.Context, string) *entities.StudentGroup
	FindFirstByName(ctx context.Context, name string) *entities.StudentGroup
	ExternalUpdate(context.Context, uuid.UUID, entities.StudentGroup) (entities.StudentGroup, error)
}

func NewStudentGroupRepository(db *sqlx.DB) StudentGroupRepository {
	config := platform.NewRepositoryConfig(
		"StudentGroupRepository",
		entities.StudentGroup{}.TableName(),
		"student group",
		[]string{"id", "slug", "name"},
		[]string{},
		[]string{"created_at", "updated_at"},
	)

	sr := &studentGroupRepository{
		BaseRepository: platform.NewBaseRepository[entities.StudentGroup](config, db),
		db:             db,
	}
	sr.logger = sr.GetLogger()

	return sr
}

type studentGroupRepository struct {
	platform.BaseRepository[entities.StudentGroup]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *studentGroupRepository) FindBySlug(ctx context.Context, slug string) *entities.StudentGroup {
	return r.FindFirstBy(ctx, "slug", slug)
}
func (r *studentGroupRepository) FindFirstByName(ctx context.Context, name string) *entities.StudentGroup {
	return r.FindFirstBy(ctx, "name", name)
}
func (r *studentGroupRepository) ExternalUpdate(
	ctx context.Context, id uuid.UUID, group entities.StudentGroup,
) (entities.StudentGroup, error) {
	return r.UpdateFields(ctx, id, []string{"slug", "name"}, group)
}
