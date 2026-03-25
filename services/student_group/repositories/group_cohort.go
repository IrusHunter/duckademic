package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/student_group/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type GroupCohortRepository interface {
	platform.BaseRepository[entities.GroupCohort]
	FindBySlug(context.Context, string) *entities.GroupCohort
	FindFirstByName(ctx context.Context, name string) *entities.GroupCohort
}

func NewGroupCohortRepository(db *sqlx.DB) GroupCohortRepository {
	config := platform.NewRepositoryConfig(
		"GroupCohortRepository",
		entities.GroupCohort{}.TableName(),
		"group_cohort",
		[]string{"id", "slug", "name", "semester_id"},
		[]string{"name", "semester_id"},
		[]string{"created_at", "updated_at"},
	)

	gr := &groupCohortRepository{
		BaseRepository: platform.NewBaseRepository[entities.GroupCohort](config, db),
		db:             db,
	}
	gr.logger = gr.GetLogger()

	return gr
}

type groupCohortRepository struct {
	platform.BaseRepository[entities.GroupCohort]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *groupCohortRepository) FindBySlug(ctx context.Context, slug string) *entities.GroupCohort {
	return r.FindFirstBy(ctx, "slug", slug)
}
func (r *groupCohortRepository) FindFirstByName(ctx context.Context, name string) *entities.GroupCohort {
	return r.FindFirstBy(ctx, "name", name)
}
