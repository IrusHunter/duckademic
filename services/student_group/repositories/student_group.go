package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/student_group/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type StudentGroupRepository interface {
	platform.BaseRepository[entities.StudentGroup]
	FindBySlug(ctx context.Context, slug string) *entities.StudentGroup
	FindFirstByName(ctx context.Context, name string) *entities.StudentGroup
}

func NewStudentGroupRepository(db *sqlx.DB) StudentGroupRepository {
	config := platform.NewRepositoryConfig(
		"StudentGroupRepository",
		entities.StudentGroup{}.TableName(),
		"student_group",
		[]string{"id", "slug", "name", "group_cohort_id"},
		[]string{"name", "group_cohort_id"},
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
