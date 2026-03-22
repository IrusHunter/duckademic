package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/IrusHunter/duckademic/services/student/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type StudentRepository interface {
	platform.BaseRepository[entities.Student]
	FindBySlug(context.Context, string) *entities.Student
	FindFirstByName(ctx context.Context, first, last string) *entities.Student
}

func NewStudentRepository(db *sqlx.DB) StudentRepository {
	config := platform.NewRepositoryConfig(
		"StudentRepository",
		entities.Student{}.TableName(),
		"student",
		[]string{"id", "slug", "first_name", "last_name", "middle_name", "phone_number", "email"},
		[]string{"first_name", "last_name", "middle_name", "phone_number", "email"},
		[]string{"created_at", "updated_at"},
	)

	sr := &studentRepository{
		BaseRepository: platform.NewBaseRepository[entities.Student](config, db),
		db:             db,
	}
	sr.logger = sr.GetLogger()

	return sr
}

type studentRepository struct {
	platform.BaseRepository[entities.Student]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *studentRepository) FindBySlug(ctx context.Context, slug string) *entities.Student {
	return r.FindFirstBy(ctx, "slug", slug)
}

func (r *studentRepository) FindFirstByName(ctx context.Context, first, last string) *entities.Student {
	query := fmt.Sprintf(
		`SELECT * FROM %s
		WHERE first_name=$1 AND last_name=$2 LIMIT 1`,
		entities.Student{}.TableName(),
	)

	var student entities.Student
	if err := r.db.GetContext(ctx, &student, query, first, last); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			r.logger.Log(contextutil.GetTraceID(ctx), "FindFirstByName",
				fmt.Sprintf("student with name %s %s not found", first, last),
				logger.RepositoryOperationSuccess,
			)
			return nil
		}
		r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "FindFirstByName",
			fmt.Errorf("failed to scan database row for %s %s: %w", first, last, err),
			logger.RepositoryScanFailed,
		)
		return nil
	}

	r.logger.Log(contextutil.GetTraceID(ctx), "FindFirstByName",
		fmt.Sprintf("%s found", student.String()),
		logger.RepositoryOperationSuccess)

	return &student
}
