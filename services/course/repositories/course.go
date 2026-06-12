package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CourseRepository interface {
	platform.BaseRepository[entities.Course]
	FindFirstByName(ctx context.Context, name string) *entities.Course
	ExternalUpdate(context.Context, uuid.UUID, entities.Course) (entities.Course, error)
	GetFullCourses(context.Context, uuid.UUID) ([]entities.Course, error)
}

func NewCourseRepository(db *sqlx.DB) CourseRepository {
	config := platform.NewRepositoryConfig(
		"CourseRepository",
		entities.Course{}.TableName(),
		entities.Course{}.EntityName(),
		[]string{"id", "slug", "name"},
		[]string{"manager_id", "description"},
		[]string{"created_at", "updated_at"},
	)

	cr := &courseRepository{
		BaseRepository: platform.NewBaseRepository[entities.Course](config, db),
		db:             db,
	}
	cr.logger = cr.GetLogger()

	return cr
}

type courseRepository struct {
	platform.BaseRepository[entities.Course]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *courseRepository) FindFirstByName(ctx context.Context, name string) *entities.Course {
	return r.FindFirstBy(ctx, "name", name)
}

func (r *courseRepository) ExternalUpdate(
	ctx context.Context, id uuid.UUID, course entities.Course,
) (entities.Course, error) {
	return r.UpdateFields(ctx, id, []string{"slug", "name"}, course)
}

type CourseFlat struct {
	ID          uuid.UUID  `db:"id"`
	ManagerID   *uuid.UUID `db:"manager_id"`
	Slug        string     `db:"slug"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`

	TeacherID        *uuid.UUID `db:"teacher_id"`
	TeacherSlug      *string    `db:"teacher_slug"`
	TeacherName      *string    `db:"teacher_name"`
	TeacherCreatedAt *time.Time `db:"teacher_created_at"`
	TeacherUpdatedAt *time.Time `db:"teacher_updated_at"`
}

func (f *CourseFlat) Convert() entities.Course {
	course := entities.Course{
		ID:          f.ID,
		ManagerID:   f.ManagerID,
		Slug:        f.Slug,
		Name:        f.Name,
		Description: f.Description,
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	}

	if f.TeacherID != nil {
		course.Manager = &entities.Teacher{
			ID:        *f.TeacherID,
			Slug:      valueOrEmpty(f.TeacherSlug),
			Name:      valueOrEmpty(f.TeacherName),
			CreatedAt: valueOrZeroTime(f.TeacherCreatedAt),
			UpdatedAt: valueOrZeroTime(f.TeacherUpdatedAt),
		}
	}

	return course
}

func valueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func valueOrZeroTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func (r *courseRepository) GetFullCourses(ctx context.Context, userID uuid.UUID) ([]entities.Course, error) {
	query := fmt.Sprintf(`
		SELECT
			c.id,
			c.manager_id,
			c.slug,
			c.name,
			c.description,
			c.created_at,
			c.updated_at,

			t.id         AS teacher_id,
			t.slug       AS teacher_slug,
			t.name       AS teacher_name,
			t.created_at AS teacher_created_at,
			t.updated_at AS teacher_updated_at

		FROM %s sc
		LEFT JOIN %s c ON sc.course_id = c.id
		LEFT JOIN %s t ON c.manager_id = t.id

		WHERE sc.student_id = $1

		ORDER BY c.name;
	`,
		entities.StudentCourse{}.TableName(),
		entities.Course{}.TableName(),
		entities.Teacher{}.TableName(),
	)

	var flats []CourseFlat
	if err := r.db.SelectContext(ctx, &flats, query, userID); err != nil {
		return nil, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetFullCourses",
			err,
			logger.RepositoryScanFailed,
		)
	}

	result := make([]entities.Course, 0, len(flats))
	for i := range flats {
		result = append(result, flats[i].Convert())
	}

	return result, nil
}
