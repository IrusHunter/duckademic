package resthandlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IrusHunter/duckademic/services/curriculum/services"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
)

type DatabaseHandler interface {
	Seed(context.Context, http.ResponseWriter, *http.Request)
	Clear(context.Context, http.ResponseWriter, *http.Request)
}

func NewDatabaseHandler(
	cs services.CurriculumService,
	ss services.SemesterService,
) DatabaseHandler {
	return &databaseHandler{
		curriculumService: cs,
		semesterService:   ss,
	}
}

type databaseHandler struct {
	curriculumService services.CurriculumService
	semesterService   services.SemesterService
}

func (h *databaseHandler) Seed(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.curriculumService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed curriculums: %w", err))
		return
	}
	if err := h.semesterService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed semesters: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
func (h *databaseHandler) Clear(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.semesterService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear semesters: %w", err))
		return
	}
	if err := h.curriculumService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear curriculums: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
