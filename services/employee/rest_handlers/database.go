package resthandlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IrusHunter/duckademic/services/employees/services"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
)

// DatabaseHandler represents a handler responsible for database-related HTTP operations.
type DatabaseHandler interface {
	// Performs database seeding operations, initializing required data.
	Seed(context.Context, http.ResponseWriter, *http.Request)
}

// NewDatabaseHandler creates a new DatabaseHandler instance.
//
// It requires the academic rank (ars), academic degree (ads), and employee (es) service.
func NewDatabaseHandler(
	ars services.AcademicRankService,
	ads services.AcademicDegreeService,
	es services.EmployeeService,
) DatabaseHandler {
	return &databaseHandler{
		academicRankService:   ars,
		academicDegreeService: ads,
		employeeService:       es,
	}
}

type databaseHandler struct {
	academicRankService   services.AcademicRankService
	academicDegreeService services.AcademicDegreeService
	employeeService       services.EmployeeService
}

func (h *databaseHandler) Seed(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.academicRankService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed academic ranks: %w", err))
		return
	}
	if err := h.academicDegreeService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed academic degrees: %w", err))
		return
	}
	if err := h.employeeService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed employees: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
