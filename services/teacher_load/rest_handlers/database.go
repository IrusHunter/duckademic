package resthandlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/IrusHunter/duckademic/services/teacher_load/services"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
)

// DatabaseHandler represents a handler responsible for database-related HTTP operations.
type DatabaseHandler interface {
	Seed(context.Context, http.ResponseWriter, *http.Request)
	Clear(context.Context, http.ResponseWriter, *http.Request)
}

func NewDatabaseHandler(
	tls services.TeacherLoadService,
	ts services.TeacherService,
	ds services.DisciplineService,
	ls services.LessonTypeService,
) DatabaseHandler {
	return &databaseHandler{
		teacherLoadService: tls,
		teacherService:     ts,
		disciplineService:  ds,
		lessonTypeService:  ls,
	}
}

type databaseHandler struct {
	teacherLoadService services.TeacherLoadService
	teacherService     services.TeacherService
	disciplineService  services.DisciplineService
	lessonTypeService  services.LessonTypeService
}

func (h *databaseHandler) Seed(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(events.ExternalSeedCooldown * 2)
		ctx := contextutil.SetTraceID(context.Background())
		h.teacherLoadService.Seed(ctx)
	}()

	jsonutil.ResponseWithJSON(w, 204, nil)
}

func (h *databaseHandler) Clear(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.teacherLoadService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear teacher loads: %w", err))
		return
	}
	if err := h.lessonTypeService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear lesson types: %w", err))
		return
	}
	if err := h.disciplineService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear disciplines: %w", err))
		return
	}
	if err := h.teacherService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear teachers: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
