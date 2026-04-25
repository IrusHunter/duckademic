package repositories

import (
	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type TeacherCourseRepository interface {
	platform.BaseRepository[entities.TeacherCourse]
}

func NewTeacherCourseRepository(db *sqlx.DB) TeacherCourseRepository {
	config := platform.NewRepositoryConfig(
		"TeacherCourseRepository",
		entities.TeacherCourse{}.TableName(),
		entities.TeacherCourse{}.EntityName(),
		[]string{"id", "course_id", "teacher_id"},
		[]string{},
		[]string{"created_at", "updated_at"},
	)

	tcr := &teacherCourseRepository{
		BaseRepository: platform.NewBaseRepository[entities.TeacherCourse](config, db),
		db:             db,
	}
	tcr.logger = tcr.GetLogger()

	return tcr
}

type teacherCourseRepository struct {
	platform.BaseRepository[entities.TeacherCourse]
	db     *sqlx.DB
	logger logger.Logger
}
