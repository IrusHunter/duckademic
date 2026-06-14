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

type TaskCommentRepository interface {
	platform.BaseRepository[entities.TaskComment]
	GetForTask(ctx context.Context, taskID uuid.UUID, isPrivate bool, studentID *uuid.UUID) ([]entities.TaskComment, error)
}

func NewTaskCommentRepository(db *sqlx.DB) TaskCommentRepository {
	config := platform.NewRepositoryConfig(
		"TaskCommentRepository",
		entities.TaskComment{}.TableName(),
		entities.TaskComment{}.EntityName(),
		[]string{"id", "task_id", "author_id", "author_name", "body", "is_private", "student_id"},
		[]string{"body"},
		[]string{"created_at", "updated_at"},
	)

	r := &taskCommentRepository{
		BaseRepository: platform.NewBaseRepository[entities.TaskComment](config, db),
		db:             db,
	}
	r.logger = r.GetLogger()
	return r
}

type taskCommentRepository struct {
	platform.BaseRepository[entities.TaskComment]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *taskCommentRepository) GetForTask(
	ctx context.Context,
	taskID uuid.UUID,
	isPrivate bool,
	studentID *uuid.UUID,
) ([]entities.TaskComment, error) {
	var query string
	var args []any

	if isPrivate && studentID != nil {
		query = fmt.Sprintf(`
			SELECT id, task_id, author_id, author_name, body, is_private, student_id, created_at, updated_at
			FROM %s
			WHERE task_id = $1 AND is_private = TRUE AND student_id = $2
			ORDER BY created_at ASC
		`, entities.TaskComment{}.TableName())
		args = []any{taskID, *studentID}
	} else {
		query = fmt.Sprintf(`
			SELECT id, task_id, author_id, author_name, body, is_private, student_id, created_at, updated_at
			FROM %s
			WHERE task_id = $1 AND is_private = FALSE
			ORDER BY created_at ASC
		`, entities.TaskComment{}.TableName())
		args = []any{taskID}
	}

	var comments []entities.TaskComment
	if err := r.db.SelectContext(ctx, &comments, query, args...); err != nil {
		return nil, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx), "GetForTask", err, logger.RepositoryScanFailed,
		)
	}
	return comments, nil
}
