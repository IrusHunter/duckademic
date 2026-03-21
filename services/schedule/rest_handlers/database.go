package resthandlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
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
	ts services.TeacherService,
) DatabaseHandler {
	return &databaseHandler{
		academicRankService: ars,
		teacherService:      ts,
	}
}

type databaseHandler struct {
	academicRankService services.AcademicRankService
	teacherService      services.TeacherService
}

func (h *databaseHandler) Seed(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(events.ExternalSeedCooldown)
		ctx := contextutil.SetTraceID(context.Background())
		h.academicRankService.Seed(ctx)
		ctx = contextutil.SetTraceID(context.Background())
		h.teacherService.Seed(ctx)

	}()

	jsonutil.ResponseWithJSON(w, 204, nil)
}
func (h *databaseHandler) Clear(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.academicRankService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear academic ranks: %w", err))
		return
	}
	if err := h.teacherService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear teachers: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
