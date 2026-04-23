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

type GroupMemberRepository interface {
	platform.BaseRepository[entities.GroupMember]
	ExternalUpdate(context.Context, uuid.UUID, entities.GroupMember) (entities.GroupMember, error)
	GetByGroupID(context.Context, uuid.UUID) ([]uuid.UUID, error)
	GetByStudentIDs(context.Context, []uuid.UUID) ([]entities.GroupMember, error)
	GetByStudentID(ctx context.Context, studentID uuid.UUID) ([]uuid.UUID, error)
}

func NewGroupMemberRepository(db *sqlx.DB) GroupMemberRepository {
	config := platform.NewRepositoryConfig(
		"GroupMemberRepository",
		entities.GroupMember{}.TableName(),
		"group member",
		[]string{"id", "student_id", "student_group_id"},
		[]string{},
		[]string{"created_at", "updated_at"},
	)

	gr := &groupMemberRepository{
		BaseRepository: platform.NewBaseRepository[entities.GroupMember](config, db),
		db:             db,
	}
	gr.logger = gr.GetLogger()

	return gr
}

type groupMemberRepository struct {
	platform.BaseRepository[entities.GroupMember]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *groupMemberRepository) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	member entities.GroupMember,
) (entities.GroupMember, error) {
	return r.UpdateFields(ctx, id, []string{"student_id", "student_group_id"}, member)
}

func (r *groupMemberRepository) GetByGroupID(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error) {
	query := fmt.Sprintf(`
		SELECT student_id
		FROM %s
		WHERE student_group_id = $1;
	`, entities.GroupMember{}.TableName())

	var ids []uuid.UUID

	if err := r.db.SelectContext(ctx, &ids, query, groupID); err != nil {
		return nil, r.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetByGroupID",
			err,
			logger.RepositoryScanFailed,
		)
	}

	return ids, nil
}
func (r *groupMemberRepository) GetByStudentIDs(ctx context.Context, studentIDs []uuid.UUID) ([]entities.GroupMember, error) {
	if len(studentIDs) == 0 {
		return []entities.GroupMember{}, nil
	}

	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT id, student_id, student_group_id
		FROM %s
		WHERE student_id IN (?);
	`, entities.GroupMember{}.TableName()), studentIDs)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)

	var members []entities.GroupMember
	if err := r.db.SelectContext(ctx, &members, query, args...); err != nil {
		return nil, r.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetByStudentIDs",
			err,
			logger.RepositoryScanFailed,
		)
	}

	return members, nil
}
func (r *groupMemberRepository) GetByStudentID(ctx context.Context, studentID uuid.UUID) ([]uuid.UUID, error) {
	query := fmt.Sprintf(`
		SELECT student_group_id
		FROM %s
		WHERE student_id = $1 AND student_group_id IS NOT NULL;
	`, entities.GroupMember{}.TableName())

	var ids []uuid.UUID

	if err := r.db.SelectContext(ctx, &ids, query, studentID); err != nil {
		return nil, r.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetByStudentID",
			err,
			logger.RepositoryScanFailed,
		)
	}

	return ids, nil
}
