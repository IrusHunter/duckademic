package repositories

import (
	"github.com/IrusHunter/duckademic/services/student_group/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type GroupCohortAssignmentRepository interface {
	platform.BaseRepository[entities.GroupCohortAssignment]
}

func NewGroupCohortAssignmentRepository(db *sqlx.DB) GroupCohortAssignmentRepository {
	config := platform.NewRepositoryConfig(
		"GroupCohortAssignmentRepository",
		entities.GroupCohortAssignment{}.TableName(),
		"group_cohort_assignment",
		[]string{"id", "group_cohort_id", "discipline_id", "lesson_type_id"},
		[]string{"group_cohort_id", "discipline_id", "lesson_type_id"},
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
