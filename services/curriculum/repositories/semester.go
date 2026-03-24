package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/curriculum/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SemesterRepository interface {
	platform.BaseRepository[entities.Semester]
	FindBySlug(context.Context, string) *entities.Semester
	FindByCurriculumIDAndNumber(context.Context, uuid.UUID, int) *entities.Semester
}

func NewSemesterRepository(db *sqlx.DB) SemesterRepository {
	config := platform.NewRepositoryConfig(
		"SemesterRepository",
		entities.Semester{}.TableName(),
		"semester",
		[]string{"id", "slug", "curriculum_id", "number"},
		[]string{"curriculum_id", "number"},
		[]string{"created_at", "updated_at"},
	)

	sr := &semesterRepository{
		BaseRepository: platform.NewBaseRepository[entities.Semester](config, db),
		db:             db,
	}
	sr.logger = sr.GetLogger()

	return sr
}

type semesterRepository struct {
	platform.BaseRepository[entities.Semester]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *semesterRepository) FindBySlug(ctx context.Context, slug string) *entities.Semester {
	return r.FindFirstBy(ctx, "slug", slug)
}
func (r *semesterRepository) FindByCurriculumIDAndNumber(
	ctx context.Context, curriculumID uuid.UUID, number int,
) *entities.Semester {
	conditions := map[string]any{
		"curriculum_id": curriculumID,
		"number":        number,
	}
	return r.FindFirstByConditions(ctx, conditions)
}
