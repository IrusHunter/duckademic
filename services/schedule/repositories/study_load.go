package repositories

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type StudyLoadRepository interface {
	platform.BaseRepository[entities.StudyLoad]
}

func NewStudyLoadRepository(db *sqlx.DB) StudyLoadRepository {
	config := platform.NewRepositoryConfig(
		"StudyLoadRepository",
		entities.StudyLoad{}.TableName(),
		"study load",
		[]string{
			"id",
			"teacher_id",
			"student_group_id",
			"discipline_id",
			"lesson_type_id",
		},
		[]string{},
		[]string{"created_at", "updated_at"},
	)

	return &studyLoadRepository{
		BaseRepository: platform.NewBaseRepository[entities.StudyLoad](config, db),
		db:             db,
	}
}

type studyLoadRepository struct {
	platform.BaseRepository[entities.StudyLoad]
	db *sqlx.DB
}
