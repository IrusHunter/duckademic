package repositories

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type TeacherSlotPriorityRepository interface {
	platform.BaseRepository[entities.TeacherSlotPriority]
}

func NewTeacherSlotPriorityRepository(db *sqlx.DB) TeacherSlotPriorityRepository {
	config := platform.NewRepositoryConfig(
		"TeacherSlotPriorityRepository",
		entities.TeacherSlotPriority{}.TableName(),
		entities.TeacherSlotPriority{}.EntityName(),
		[]string{
			"id",
			"teacher_id",
			"time_slot_id",
			"priority",
			"created_at",
			"updated_at",
		},
		[]string{"priority"},
		[]string{"created_at", "updated_at"},
	)

	return &teacherSlotPriorityRepository{
		BaseRepository: platform.NewBaseRepository[entities.TeacherSlotPriority](config, db),
		db:             db,
	}
}

type teacherSlotPriorityRepository struct {
	platform.BaseRepository[entities.TeacherSlotPriority]
	db *sqlx.DB
}
