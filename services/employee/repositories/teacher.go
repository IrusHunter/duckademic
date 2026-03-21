package repositories

import (
	"github.com/IrusHunter/duckademic/services/employees/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type TeacherRepository interface {
	platform.BaseRepository[entities.Teacher]
}

func NewTeacherRepository(db *sqlx.DB) TeacherRepository {
	config := platform.NewRepositoryConfig("TeacherRepository", "teachers", "teacher",
		[]string{"employee_id", "email", "academic_degree_id", "academic_rank_id"},
		[]string{"email", "academic_degree_id", "academic_rank_id"},
		[]string{"created_at", "updated_at"},
	)

	ts := &teacherRepository{
		BaseRepository: platform.NewBaseRepository[entities.Teacher](config, db),
		db:             db,
	}
	ts.logger = ts.GetLogger()
	return ts
}

type teacherRepository struct {
	platform.BaseRepository[entities.Teacher]
	db     *sqlx.DB
	logger logger.Logger
}
