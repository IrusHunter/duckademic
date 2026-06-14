package repositories

import (
	"context"
	"fmt"

	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type NotificationRepository interface {
	Create(ctx context.Context, n entities.Notification) error
	GetForRecipient(ctx context.Context, recipientID uuid.UUID) ([]entities.Notification, error)
	MarkAllReadForRecipient(ctx context.Context, recipientID uuid.UUID) error
}

func NewNotificationRepository(db *sqlx.DB) NotificationRepository {
	return &notificationRepository{
		db:     db,
		logger: logger.NewLogger("NotificationRepository.txt", "NotificationRepository"),
	}
}

type notificationRepository struct {
	db     *sqlx.DB
	logger logger.Logger
}

func (r *notificationRepository) Create(ctx context.Context, n entities.Notification) error {
	query := `
		INSERT INTO notifications (id, recipient_id, type, title, task_id, course_id, task_student_id, is_read, created_at)
		VALUES (:id, :recipient_id, :type, :title, :task_id, :course_id, :task_student_id, :is_read, NOW())
	`
	if _, err := r.db.NamedExecContext(ctx, query, n); err != nil {
		return r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Create",
			fmt.Errorf("failed to create notification: %w", err), logger.RepositoryQueryFailed)
	}
	return nil
}

func (r *notificationRepository) GetForRecipient(ctx context.Context, recipientID uuid.UUID) ([]entities.Notification, error) {
	query := `
		SELECT id, recipient_id, type, title, task_id, course_id, task_student_id, is_read, created_at
		FROM notifications
		WHERE recipient_id = $1
		ORDER BY created_at DESC
		LIMIT 50
	`
	var notifs []entities.Notification
	if err := r.db.SelectContext(ctx, &notifs, query, recipientID); err != nil {
		return nil, r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "GetForRecipient",
			fmt.Errorf("failed to get notifications: %w", err), logger.RepositoryScanFailed)
	}
	return notifs, nil
}

func (r *notificationRepository) MarkAllReadForRecipient(ctx context.Context, recipientID uuid.UUID) error {
	query := `UPDATE notifications SET is_read = true WHERE recipient_id = $1 AND is_read = false`
	if _, err := r.db.ExecContext(ctx, query, recipientID); err != nil {
		return r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "MarkAllReadForRecipient",
			fmt.Errorf("failed to mark notifications read: %w", err), logger.RepositoryQueryFailed)
	}
	return nil
}
