package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GroupMemberRepository interface {
	platform.BaseRepository[entities.GroupMember]
	ExternalUpdate(ctx context.Context, id uuid.UUID, member entities.GroupMember) (entities.GroupMember, error)
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
