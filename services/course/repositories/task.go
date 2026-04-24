package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type TaskRepository interface {
	platform.BaseRepository[entities.Task]
	FindFirstByTitle(context.Context, string) *entities.Task
}

func NewTaskRepository(db *sqlx.DB) TaskRepository {
	config := platform.NewRepositoryConfig(
		"TaskRepository",
		entities.Task{}.TableName(),
		entities.Task{}.EntityName(),
		[]string{"id", "course_id", "slug", "title", "description", "max_mark", "deadline"},
		[]string{"title", "description", "max_mark", "deadline"},
		[]string{"created_at", "updated_at"},
	)

	tr := &taskRepository{
		BaseRepository: platform.NewBaseRepository[entities.Task](config, db),
		db:             db,
	}
	tr.logger = tr.GetLogger()

	return tr
}

type taskRepository struct {
	platform.BaseRepository[entities.Task]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *taskRepository) FindFirstByTitle(ctx context.Context, title string) *entities.Task {
	return r.FindFirstBy(ctx, "title", title)
}
