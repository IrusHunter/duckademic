package repositories

import (
	"context"
	"fmt"

	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type StudentCourseRepository interface {
	platform.BaseRepository[entities.StudentCourse]
	GetCoursesForStudent(ctx context.Context, studentID uuid.UUID) ([]entities.Course, error)
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

	return scr
}

type studentCourseRepository struct {
	platform.BaseRepository[entities.StudentCourse]
	db *sqlx.DB
}

func (r *studentCourseRepository) GetCoursesForStudent(ctx context.Context, studentID uuid.UUID) ([]entities.Course, error) {
	query := fmt.Sprintf(`
		SELECT 
			c.id,
			c.manager_id,
			c.slug,
			c.name,
			c.description,
			c.created_at,
			c.updated_at
		FROM %s sc
		JOIN %s c ON sc.course_id = c.id
		WHERE sc.student_id = ?
		ORDER BY c.created_at DESC;
	`,
		entities.StudentCourse{}.TableName(),
		entities.Course{}.TableName(),
	)

	query = r.db.Rebind(query)

	var courses []entities.Course
	if err := r.db.SelectContext(ctx, &courses, query, studentID); err != nil {
		return nil, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetCoursesForStudent",
			err,
			logger.RepositoryScanFailed,
		)
	}

	return courses, nil
}
