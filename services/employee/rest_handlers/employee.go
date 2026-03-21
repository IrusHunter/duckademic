package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/employee/entities"
	"github.com/IrusHunter/duckademic/services/employee/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// EmployeeHandler represents a handler responsible for Employee-related HTTP operations.
type EmployeeHandler interface {
	platform.BaseHandler[entities.Employee]
}

// NewEmployeeHandler creates a new EmployeeHandler instance.
//
// It requires an employee services (es).
func NewEmployeeHandler(es services.EmployeeService) EmployeeHandler {
	hc := platform.NewHandlerConfig("EmployeeHandler", "employee")

	return &employeeHandler{
		BaseHandler: platform.NewBaseHandler(hc, es),
		service:     es,
	}
}

type employeeHandler struct {
	platform.BaseHandler[entities.Employee]
	service services.EmployeeService
}
