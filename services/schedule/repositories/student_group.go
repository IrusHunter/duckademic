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

type StudentGroupRepository interface {
	platform.BaseRepository[entities.StudentGroup]
	FindBySlug(context.Context, string) *entities.StudentGroup
	FindFirstByName(ctx context.Context, name string) *entities.StudentGroup
	ExternalUpdate(context.Context, uuid.UUID, entities.StudentGroup) (entities.StudentGroup, error)
	GetByGroupCohortID(context.Context, uuid.UUID) ([]entities.StudentGroup, error)
}

func NewStudentGroupRepository(db *sqlx.DB) StudentGroupRepository {
	config := platform.NewRepositoryConfig(
		"StudentGroupRepository",
		entities.StudentGroup{}.TableName(),
		"student group",
		[]string{"id", "slug", "name", "group_cohort_id"},
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
	return r.UpdateFields(ctx, id, []string{"slug", "name", "group_cohort_id"}, group)
}

func (r *studentGroupRepository) GetByGroupCohortID(ctx context.Context, cohortID uuid.UUID) ([]entities.StudentGroup, error) {
	query := fmt.Sprintf(`
		SELECT sg.id, sg.name
		FROM %s sg
		WHERE sg.group_cohort_id = $1;
	`, entities.StudentGroup{}.TableName())

	var studentGroups []entities.StudentGroup

	if err := r.db.SelectContext(ctx, &studentGroups, query, cohortID); err != nil {
		return nil, r.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetByGroupCohortID",
			err,
			logger.RepositoryScanFailed,
		)
	}

	return studentGroups, nil
}
