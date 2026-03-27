package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/teacher_load/entities"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TeacherRepository interface {
	platform.BaseRepository[entities.Teacher]
	FindByName(context.Context, string) *entities.Teacher
	ExternalUpdate(context.Context, uuid.UUID, entities.Teacher) (entities.Teacher, error)
}

func NewTeacherRepository(db *sqlx.DB) TeacherRepository {
	config := platform.NewRepositoryConfig("TeacherRepository", entities.Teacher{}.TableName(),
		"teacher", []string{"id", "name"}, []string{""}, []string{"created_at", "updated_at"},
	)
	return &teacherRepository{
		BaseRepository: platform.NewBaseRepository[entities.Teacher](config, db),
		db:             db,
	}
}

type teacherRepository struct {
	platform.BaseRepository[entities.Teacher]
	db *sqlx.DB
}

func (r *teacherRepository) FindByName(ctx context.Context, name string) *entities.Teacher {
	return r.FindFirstBy(ctx, "name", name)
}
func (r *teacherRepository) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	teacher entities.Teacher,
) (entities.Teacher, error) {
	return r.UpdateFields(ctx, id, []string{"name"}, teacher)
}
