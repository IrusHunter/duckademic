package repositories

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type LessonSlotRepository interface {
	platform.BaseRepository[entities.LessonSlot]
}

func NewLessonSlotRepository(db *sqlx.DB) LessonSlotRepository {
	config := platform.NewRepositoryConfig(
		"LessonSlotRepository",
		entities.LessonSlot{}.TableName(),
		"lesson slot",
		[]string{
			"id",
			"slot",
			"weekday",
			"start_time",
			"duration",
		},
		[]string{},
		[]string{"created_at", "updated_at"},
	)

	return &lessonSlotRepository{
		BaseRepository: platform.NewBaseRepository[entities.LessonSlot](config, db),
		db:             db,
	}
}

type lessonSlotRepository struct {
	platform.BaseRepository[entities.LessonSlot]
	db *sqlx.DB
}
