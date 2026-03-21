package resthandlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IrusHunter/duckademic/services/employee/services"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
)

// DatabaseHandler represents a handler responsible for database-related HTTP operations.
type DatabaseHandler interface {
	// Performs database seeding operations, initializing required data.
	Seed(context.Context, http.ResponseWriter, *http.Request)
	Clear(context.Context, http.ResponseWriter, *http.Request)
}

func NewDatabaseHandler(
	ars services.AcademicRankService,
	ads services.AcademicDegreeService,
	es services.EmployeeService,
	ts services.TeacherService,
) DatabaseHandler {
	return &databaseHandler{
		academicRankService:   ars,
		academicDegreeService: ads,
		employeeService:       es,
		teacherService:        ts,
	}
}

type databaseHandler struct {
	academicRankService   services.AcademicRankService
	academicDegreeService services.AcademicDegreeService
	employeeService       services.EmployeeService
	teacherService        services.TeacherService
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
	if err := h.teacherService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed teachers: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
func (h *databaseHandler) Clear(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.teacherService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear teachers: %w", err))
		return
	}
	if err := h.academicRankService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear academic ranks: %w", err))
		return
	}
	if err := h.academicDegreeService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear academic degrees: %w", err))
		return
	}
	if err := h.employeeService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear employees: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
