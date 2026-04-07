package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GroupCohortAssignmentRepository interface {
	platform.BaseRepository[entities.GroupCohortAssignment]
	ExternalUpdate(ctx context.Context, id uuid.UUID, assignment entities.GroupCohortAssignment) (entities.GroupCohortAssignment, error)
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
