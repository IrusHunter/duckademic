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

type GroupCohortAssignmentRepository interface {
	platform.BaseRepository[entities.GroupCohortAssignment]
	ExternalUpdate(context.Context, uuid.UUID, entities.GroupCohortAssignment) (entities.GroupCohortAssignment, error)
	GetByGroupCohortID(context.Context, uuid.UUID) ([]entities.GroupCohortAssignment, error)
}

func NewGroupCohortAssignmentRepository(db *sqlx.DB) GroupCohortAssignmentRepository {
	config := platform.NewRepositoryConfig(
		"GroupCohortAssignmentRepository",
		entities.GroupCohortAssignment{}.TableName(),
		"group cohort assignment",
		[]string{"id", "group_cohort_id", "discipline_id", "lesson_type_id"},
		[]string{},
		[]string{"created_at", "updated_at"},
	)

	gr := &groupCohortAssignmentRepository{
		BaseRepository: platform.NewBaseRepository[entities.GroupCohortAssignment](config, db),
		db:             db,
	}
	gr.logger = gr.GetLogger()

	return gr
}

type groupCohortAssignmentRepository struct {
	platform.BaseRepository[entities.GroupCohortAssignment]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *groupCohortAssignmentRepository) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	assignment entities.GroupCohortAssignment,
) (entities.GroupCohortAssignment, error) {
	return r.UpdateFields(
		ctx,
		id,
		[]string{"group_cohort_id", "discipline_id", "lesson_type_id"},
		assignment,
	)
}

func (r *groupCohortAssignmentRepository) GetByGroupCohortID(
	ctx context.Context,
	cohortID uuid.UUID,
) ([]entities.GroupCohortAssignment, error) {
	query := fmt.Sprintf(`
		SELECT gca.id, gca.group_cohort_id, gca.discipline_id, gca.lesson_type_id
		FROM %s gca
		WHERE gca.group_cohort_id = $1;
	`, entities.GroupCohortAssignment{}.TableName())

	var groupCohortAssignments []entities.GroupCohortAssignment

	if err := r.db.SelectContext(ctx, &groupCohortAssignments, query, cohortID); err != nil {
		return nil, r.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetByGroupCohortID",
			err,
			logger.RepositoryScanFailed,
		)
	}

	return groupCohortAssignments, nil
}
