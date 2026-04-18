package resthandlers

import (
	"context"
	"net/http"

	"github.com/IrusHunter/duckademic/shared/jsonutil"
)

type DatabaseHandler interface {
	Seed(context.Context, http.ResponseWriter, *http.Request)
	Clear(context.Context, http.ResponseWriter, *http.Request)
}

func NewDatabaseHandler() DatabaseHandler {
	return &databaseHandler{}
}

type databaseHandler struct {
}

func (h *databaseHandler) Seed(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// go func() {
	// 	time.Sleep(events.ExternalSeedCooldown)
	// 	ctx = contextutil.SetTraceID(context.Background())
	// 	h.studentService.Seed(ctx)
	// }()

	jsonutil.ResponseWithJSON(w, 204, nil)
}
func (h *databaseHandler) Clear(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// if err := h.studentService.Clear(ctx); err != nil {
	// 	jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear students: %w", err))
	// 	return
	// }

	jsonutil.ResponseWithJSON(w, 204, nil)
}
