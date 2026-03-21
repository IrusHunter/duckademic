package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/IrusHunter/duckademic/services/employees/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

// EmployeeRepository represents a storage for employee entities.
type EmployeeRepository interface {
	platform.BaseRepository[entities.Employee]
	// FindBySlug returns a pointer to the employee from database with the given slug.
	FindBySlug(context.Context, string) *entities.Employee
	FindFirstByName(ctx context.Context, first, last string) *entities.Employee
}

// NewEmployeeRepository creates a new EmployeeRepository instance.
//
// It requires a database connection (db).
func NewEmployeeRepository(db *sqlx.DB) EmployeeRepository {
	config := platform.NewRepositoryConfig("EmployeeRepository", "employees", "employee",
		[]string{"id", "slug", "first_name", "last_name", "middle_name", "phone_number"},
		[]string{"first_name", "last_name", "middle_name", "phone_number"},
		[]string{"created_at", "updated_at"},
	)

	er := &employeeRepository{
		BaseRepository: platform.NewBaseRepository[entities.Employee](config, db),
		db:             db,
	}
	er.logger = er.GetLogger()

	return er
}

type employeeRepository struct {
	platform.BaseRepository[entities.Employee]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *employeeRepository) FindBySlug(ctx context.Context, slug string) *entities.Employee {
	return r.FindFirstBy(ctx, "slug", slug)
}
func (r *employeeRepository) FindByTitle(ctx context.Context, title string) *entities.Employee {
	return r.FindFirstBy(ctx, "title", title)
}
func (r *employeeRepository) FindFirstByName(ctx context.Context, first, last string) *entities.Employee {
	query := fmt.Sprintf(
		`SELECT * FROM %s
		WHERE first_name=$1 AND last_name=$2 LIMIT 1`,
		entities.Employee{}.TableName(),
	)

	var employee entities.Employee
	if err := r.db.GetContext(ctx, &employee, query, first, last); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			r.logger.Log(contextutil.GetTraceID(ctx), "FindFirstByName",
				fmt.Sprintf("employee with name %s %s not found", first, last),
				logger.RepositoryOperationSuccess,
			)
			return nil
		}
		r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "FindFirstByName",
			fmt.Errorf("failed to scan database row for %s %s: %w", first, last, err), logger.RepositoryScanFailed,
		)
		return nil
	}

	r.logger.Log(contextutil.GetTraceID(ctx), "FindFirstByName",
		fmt.Sprintf("%s found", employee.String()),
		logger.RepositoryOperationSuccess)
	return &employee
}
