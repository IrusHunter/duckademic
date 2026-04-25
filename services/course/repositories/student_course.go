package repositories

import (
	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type StudentCourseRepository interface {
	platform.BaseRepository[entities.StudentCourse]
}

func NewStudentCourseRepository(db *sqlx.DB) StudentCourseRepository {
	config := platform.NewRepositoryConfig(
		"StudentCourseRepository",
		entities.StudentCourse{}.TableName(),
		entities.StudentCourse{}.EntityName(),
		[]string{"id", "course_id", "student_id"},
		[]string{},
		[]string{"created_at", "updated_at"},
	)

	scr := &studentCourseRepository{
		BaseRepository: platform.NewBaseRepository[entities.StudentCourse](config, db),
		db:             db,
	}
	scr.logger = scr.GetLogger()

	return scr
}

type studentCourseRepository struct {
	platform.BaseRepository[entities.StudentCourse]
	db     *sqlx.DB
	logger logger.Logger
}
