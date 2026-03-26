package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type StudentRepository interface {
	platform.BaseRepository[entities.Student]
	FindBySlug(context.Context, string) *entities.Student
	FindFirstByName(ctx context.Context, name string) *entities.Student
	ExternalUpdate(context.Context, uuid.UUID, entities.Student) (entities.Student, error)
}

func NewStudentRepository(db *sqlx.DB) StudentRepository {
	config := platform.NewRepositoryConfig(
		"StudentRepository",
		entities.Student{}.TableName(),
		"student",
		[]string{"id", "slug", "name"},
		[]string{},
		[]string{"created_at", "updated_at"},
	)

	sr := &studentRepository{
		BaseRepository: platform.NewBaseRepository[entities.Student](config, db),
		db:             db,
	}
	sr.logger = sr.GetLogger()

	return sr
}

type studentRepository struct {
	platform.BaseRepository[entities.Student]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *studentRepository) FindBySlug(ctx context.Context, slug string) *entities.Student {
	return r.FindFirstBy(ctx, "slug", slug)
}
func (r *studentRepository) FindFirstByName(ctx context.Context, name string) *entities.Student {
	return r.FindFirstBy(ctx, "name", name)
}

func (r *studentRepository) ExternalUpdate(
	ctx context.Context, id uuid.UUID, student entities.Student,
) (entities.Student, error) {
	return r.UpdateFields(ctx, id, []string{"slug", "name"}, student)
}
