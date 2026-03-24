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
) DatabaseHandler {
	return &databaseHandler{
		curriculumService: cs,
	}
}

type databaseHandler struct {
	curriculumService services.CurriculumService
}

func (h *databaseHandler) Seed(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.curriculumService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed curriculums: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
func (h *databaseHandler) Clear(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.curriculumService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear curriculums: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
