package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/teacher_load/entities"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GroupCohortRepository interface {
	platform.BaseRepository[entities.GroupCohort]
	FindBySlug(context.Context, string) *entities.GroupCohort
	FindFirstByName(context.Context, string) *entities.GroupCohort
	ExternalUpdate(context.Context, uuid.UUID, entities.GroupCohort) (entities.GroupCohort, error)
}

func NewGroupCohortRepository(db *sqlx.DB) GroupCohortRepository {
	config := platform.NewRepositoryConfig(
		"GroupCohortRepository",
		entities.GroupCohort{}.TableName(),
		"group_cohort",
		[]string{"id", "slug", "name"},
		[]string{""},
		[]string{"created_at", "updated_at"},
	)

	return &groupCohortRepository{
		BaseRepository: platform.NewBaseRepository[entities.GroupCohort](config, db),
		db:             db,
	}
}

type groupCohortRepository struct {
	platform.BaseRepository[entities.GroupCohort]
	db *sqlx.DB
}

func (r *groupCohortRepository) FindBySlug(ctx context.Context, slug string) *entities.GroupCohort {
	return r.FindFirstBy(ctx, "slug", slug)
}
func (r *groupCohortRepository) FindFirstByName(ctx context.Context, name string) *entities.GroupCohort {
	return r.FindFirstBy(ctx, "name", name)
}
func (r *groupCohortRepository) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	groupCohort entities.GroupCohort,
) (entities.GroupCohort, error) {
	return r.UpdateFields(ctx, id, []string{"slug", "name"}, groupCohort)
}
