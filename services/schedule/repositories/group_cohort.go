package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GroupCohortRepository interface {
	platform.BaseRepository[entities.GroupCohort]
	ExternalUpdate(ctx context.Context, id uuid.UUID, cohort entities.GroupCohort) (entities.GroupCohort, error)
}

func NewGroupCohortRepository(db *sqlx.DB) GroupCohortRepository {
	config := platform.NewRepositoryConfig(
		"GroupCohortRepository",
		entities.GroupCohort{}.TableName(),
		"group cohort",
		[]string{"id", "slug", "name"},
		[]string{},
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

func (r *groupCohortRepository) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	cohort entities.GroupCohort,
) (entities.GroupCohort, error) {
	return r.UpdateFields(ctx, id, []string{"slug", "name"}, cohort)
}
