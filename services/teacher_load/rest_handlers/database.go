package resthandlers

import (
	"context"
	"net/http"

	"github.com/IrusHunter/duckademic/shared/jsonutil"
)

// DatabaseHandler represents a handler responsible for database-related HTTP operations.
type DatabaseHandler interface {
	// Performs database seeding operations, initializing required data.
	Seed(context.Context, http.ResponseWriter, *http.Request)
	Clear(context.Context, http.ResponseWriter, *http.Request)
}

func NewDatabaseHandler() DatabaseHandler {
	return &databaseHandler{}
}

type databaseHandler struct{}

func (h *databaseHandler) Seed(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	jsonutil.ResponseWithJSON(w, 204, nil)
}
func (h *databaseHandler) Clear(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// if err := h.groupMemberService.Clear(ctx); err != nil {
	// 	jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear group members: %w", err))
	// 	return
	// }

	jsonutil.ResponseWithJSON(w, 204, nil)
}
