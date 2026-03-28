package repositories

import (
	"github.com/IrusHunter/duckademic/services/teacher_load/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type TeacherLoadRepository interface {
	platform.BaseRepository[entities.TeacherLoad]
}

func NewTeacherLoadRepository(db *sqlx.DB) TeacherLoadRepository {
	config := platform.NewRepositoryConfig(
		"TeacherLoadRepository",
		entities.TeacherLoad{}.TableName(),
		"teacher_load",
		[]string{"id", "teacher_id", "discipline_id", "lesson_type_id", "group_cohort_id", "group_count"},
		[]string{"teacher_id", "discipline_id", "lesson_type_id", "group_cohort_id", "group_count"},
		[]string{"created_at", "updated_at"},
	)

	tr := &teacherLoadRepository{
		BaseRepository: platform.NewBaseRepository[entities.TeacherLoad](config, db),
		db:             db,
	}
	tr.logger = tr.GetLogger()

	return tr
}

type teacherLoadRepository struct {
	platform.BaseRepository[entities.TeacherLoad]
	db     *sqlx.DB
	logger logger.Logger
}
