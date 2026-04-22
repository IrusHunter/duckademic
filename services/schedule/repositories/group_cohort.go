package repositories

import (
	"context"
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GroupCohortRepository interface {
	platform.BaseRepository[entities.GroupCohort]
	ExternalUpdate(context.Context, uuid.UUID, entities.GroupCohort) (entities.GroupCohort, error)
	GetBySemesterID(context.Context, uuid.UUID) ([]entities.GroupCohort, error)
}

func NewGroupCohortRepository(db *sqlx.DB) GroupCohortRepository {
	config := platform.NewRepositoryConfig(
		"GroupCohortRepository",
		entities.GroupCohort{}.TableName(),
		"group cohort",
		[]string{"id", "slug", "name", "semester_id"},
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
	return r.UpdateFields(ctx, id, []string{"slug", "name", "semester_id"}, cohort)
}

func (r *groupCohortRepository) GetBySemesterID(ctx context.Context, semesterID uuid.UUID) ([]entities.GroupCohort, error) {
	query := fmt.Sprintf(`
		SELECT id, slug, name, semester_id
		FROM %s
		WHERE semester_id = $1;
	`, entities.GroupCohort{}.TableName())

	var cohorts []entities.GroupCohort

	if err := r.db.SelectContext(ctx, &cohorts, query, semesterID); err != nil {
		return nil, r.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetBySemesterID",
			err,
			logger.RepositoryScanFailed,
		)
	}

	return cohorts, nil
}
