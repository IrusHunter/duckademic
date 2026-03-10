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
// It requires a academic rank service.
func NewDatabaseHandler(ars services.AcademicRankService) DatabaseHandler {
	return &databaseHandler{academicRankService: ars}
}

type databaseHandler struct {
	academicRankService services.AcademicRankService
}

func (h *databaseHandler) Seed(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.academicRankService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed academic ranks: %w", err))
		return
	}
}
