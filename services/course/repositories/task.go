package repositories

import (
	"context"
	"fmt"

	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TaskRepository interface {
	platform.BaseRepository[entities.Task]
	FindFirstByTitle(context.Context, string) *entities.Task
	GetTasksByCourseID(context.Context, uuid.UUID) ([]entities.Task, error)
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
func (r *taskRepository) GetTasksByCourseID(ctx context.Context, courseID uuid.UUID) ([]entities.Task, error) {
	query := fmt.Sprintf(`
		SELECT 
			id,
			course_id,
			slug,
			title,
			description,
			max_mark,
			deadline,
			created_at,
			updated_at
		FROM %s
		WHERE course_id = ?
		ORDER BY created_at DESC;
	`, entities.Task{}.TableName())

	query = r.db.Rebind(query)

	var tasks []entities.Task
	if err := r.db.SelectContext(ctx, &tasks, query, courseID); err != nil {
		return nil, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetTasksByCourseID",
			err,
			logger.RepositoryScanFailed,
		)
	}

	return tasks, nil
}
