package repositories

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type LessonOccurrenceRepository interface {
	platform.BaseRepository[entities.LessonOccurrence]
}

func NewLessonOccurrenceRepository(db *sqlx.DB) LessonOccurrenceRepository {
	config := platform.NewRepositoryConfig(
		"LessonOccurrenceRepository",
		entities.LessonOccurrence{}.TableName(),
		entities.LessonOccurrence{}.EntityName(),
		[]string{
			"id",
			"study_load_id",
			"teacher_id",
			"student_group_id",
			"lesson_slot_id",
			"date",
			"classroom_id",
			"status",
		},
		[]string{},
		[]string{"created_at", "updated_at"},
	)

	return &lessonOccurrenceRepository{
		BaseRepository: platform.NewBaseRepository[entities.LessonOccurrence](config, db),
		db:             db,
	}
}

type lessonOccurrenceRepository struct {
	platform.BaseRepository[entities.LessonOccurrence]
	db *sqlx.DB
}
