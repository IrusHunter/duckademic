package repositories

import (
	"github.com/IrusHunter/duckademic/services/student_group/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type GroupMemberRepository interface {
	platform.BaseRepository[entities.GroupMember]
}

func NewGroupMemberRepository(db *sqlx.DB) GroupMemberRepository {
	config := platform.NewRepositoryConfig(
		"GroupMembersRepository",
		entities.GroupMember{}.TableName(),
		"group_member",
		[]string{"id", "student_id", "group_cohort_id", "student_group_id"},
		[]string{"student_id", "group_cohort_id", "student_group_id"},
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
