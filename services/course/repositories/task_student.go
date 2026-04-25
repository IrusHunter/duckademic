package repositories

import (
	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type TaskStudentRepository interface {
	platform.BaseRepository[entities.TaskStudent]
}

func NewTaskStudentRepository(db *sqlx.DB) TaskStudentRepository {
	config := platform.NewRepositoryConfig(
		"TaskStudentRepository",
		entities.TaskStudent{}.TableName(),
		entities.TaskStudent{}.EntityName(),
		[]string{"id", "task_id", "student_id", "mark", "submission_time"},
		[]string{},
		[]string{"created_at", "updated_at"},
	)

	tsr := &taskStudentRepository{
		BaseRepository: platform.NewBaseRepository[entities.TaskStudent](config, db),
		db:             db,
	}
	tsr.logger = tsr.GetLogger()

	return tsr
}

type taskStudentRepository struct {
	platform.BaseRepository[entities.TaskStudent]
	db     *sqlx.DB
	logger logger.Logger
}
