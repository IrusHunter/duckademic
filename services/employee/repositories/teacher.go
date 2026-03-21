package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/IrusHunter/duckademic/services/employee/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TeacherRepository interface {
	platform.BaseRepository[entities.Teacher]
	Fill(context.Context, uuid.UUID) *entities.Teacher
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

func (r *teacherRepository) Fill(ctx context.Context, id uuid.UUID) *entities.Teacher {
	query := `
		SELECT 
			t.*,

			e.id "employee.id",
			e.first_name "employee.first_name",
			e.last_name "employee.last_name",

			ad.id "academic_degree.id",
			ad.title "academic_degree.title",

			ar.id "academic_rank.id",
			ar.title "academic_rank.title"

		FROM teachers t
		LEFT JOIN employees e ON t.employee_id = e.id
		LEFT JOIN academic_degrees ad ON t.academic_degree_id = ad.id
		LEFT JOIN academic_ranks ar ON t.academic_rank_id = ar.id

		WHERE t.id = $1
		LIMIT 1
	`

	var teacher entities.Teacher

	if err := r.db.GetContext(ctx, &teacher, query, id); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			r.logger.Log(contextutil.GetTraceID(ctx), "Fill",
				fmt.Sprintf("teacher with id %q not found", id),
				logger.RepositoryOperationSuccess,
			)
			return nil
		}

		r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "FindFirstDetailedBy",
			fmt.Errorf("failed to scan teacher with id %q: %w", id, err),
			logger.RepositoryScanFailed,
		)
		return nil
	}

	r.logger.Log(contextutil.GetTraceID(ctx), "FindFirstDetailedBy",
		fmt.Sprintf("teacher %s found with relations", teacher.String()),
		logger.RepositoryOperationSuccess,
	)

	return &teacher
}
