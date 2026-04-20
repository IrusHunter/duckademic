package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/employee/entities"
	"github.com/IrusHunter/duckademic/services/employee/repositories"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

// EmployeeService provides operations to initialize and manage employees.
type EmployeeService interface {
	platform.BaseService[entities.Employee]
}

// NewEmployeeService creates a new EmployeeService instance.
//
// It requires an employee repository (er).
func NewEmployeeService(er repositories.EmployeeRepository) EmployeeService {
	sc := platform.NewServiceConfig("EmployeeRepository", filepath.Join("data", "employees.json"), "employee")

	res := &employeeRepository{
		repository: er,
	}
	res.BaseService = platform.NewBaseService(sc, er,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Employee]{
			platform.OnAddPrepare:   res.onAddPrepare,
			platform.ValidateEntity: res.validateEntity,
		},
	)

	return res
}

type employeeRepository struct {
	platform.BaseService[entities.Employee]
	repository repositories.EmployeeRepository
}

func (s *employeeRepository) validateEntity(ctx context.Context, employee *entities.Employee) error {
	if err := employee.ValidateFirstName(); err != nil {
		return err
	}
	if err := employee.ValidateLastName(); err != nil {
		return err
	}

	return nil
}
func (s *employeeRepository) onAddPrepare(ctx context.Context, employee *entities.Employee) error {
	slug := slug.Make(employee.GetFullName())
	if other := s.repository.FindBySlug(ctx, slug); other != nil {
		return fmt.Errorf("employee with slug %q already exists", slug)
	}
	employee.ID = uuid.New()
	employee.Slug = slug

	return nil
}
