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
	ds services.DisciplineService,
	lts services.LessonTypeService,
	ltas services.LessonTypeAssignmentService,
) DatabaseHandler {
	return &databaseHandler{
		academicRankService:         ars,
		teacherService:              ts,
		disciplineService:           ds,
		lessonTypeService:           lts,
		lessonTypeAssignmentService: ltas,
	}
}

type databaseHandler struct {
	academicRankService         services.AcademicRankService
	teacherService              services.TeacherService
	disciplineService           services.DisciplineService
	lessonTypeService           services.LessonTypeService
	lessonTypeAssignmentService services.LessonTypeAssignmentService
}

func (h *databaseHandler) Seed(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(events.ExternalSeedCooldown)
		ctx := contextutil.SetTraceID(context.Background())
		h.academicRankService.Seed(ctx)
		ctx = contextutil.SetTraceID(context.Background())
		h.teacherService.Seed(ctx)
		ctx = contextutil.SetTraceID(context.Background())
		h.disciplineService.Seed(ctx)
		ctx = contextutil.SetTraceID(context.Background())
		h.lessonTypeService.Seed(ctx)
		ctx = contextutil.SetTraceID(context.Background())
		h.lessonTypeAssignmentService.Seed(ctx)
	}()

	jsonutil.ResponseWithJSON(w, 204, nil)
}
func (h *databaseHandler) Clear(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.lessonTypeAssignmentService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear lesson type assignments: %w", err))
		return
	}
	if err := h.academicRankService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear academic ranks: %w", err))
		return
	}
	if err := h.teacherService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear teachers: %w", err))
		return
	}
	if err := h.disciplineService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear disciplines: %w", err))
		return
	}
	if err := h.lessonTypeService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear lesson types: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
