package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/employees/entities"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

// EmployeeRepository represents a storage for employee entities.
type EmployeeRepository interface {
	platform.BaseRepository[entities.Employee]
	// FindBySlug returns a pointer to the employee from database with the given slug.
	FindBySlug(context.Context, string) *entities.Employee
}

// NewEmployeeRepository creates a new EmployeeRepository instance.
//
// It requires a database connection (db).
func NewEmployeeRepository(db *sqlx.DB) EmployeeRepository {
	config := platform.NewRepositoryConfig("EmployeeRepository", "employees", "employee",
		[]string{"id", "slug", "first_name", "last_name", "middle_name", "phone_number"},
		[]string{"id", "slug", "first_name", "last_name", "middle_name", "phone_number",
			"created_at", "updated_at", "deleted_at"},
		[]string{"first_name", "last_name", "middle_name", "phone_number"},
		[]string{"created_at", "updated_at"},
	)
	return &employeeRepository{
		BaseRepository: platform.NewBaseRepository[entities.Employee](config, db),
		db:             db,
	}
}

type employeeRepository struct {
	platform.BaseRepository[entities.Employee]
	db *sqlx.DB
}

func (r *employeeRepository) FindBySlug(ctx context.Context, slug string) *entities.Employee {
	return r.FindFirstBy(ctx, "slug", slug)
}
