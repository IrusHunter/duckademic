package resthandlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IrusHunter/duckademic/services/course/services"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
)

type DatabaseHandler interface {
	Seed(context.Context, http.ResponseWriter, *http.Request)
	Clear(context.Context, http.ResponseWriter, *http.Request)
}

func NewDatabaseHandler(
	ss services.StudentService,
	ts services.TeacherService,
) DatabaseHandler {
	return &databaseHandler{
		studentService: ss,
		teacherService: ts,
	}
}

type databaseHandler struct {
	studentService services.StudentService
	teacherService services.TeacherService
}

func (h *databaseHandler) Seed(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	jsonutil.ResponseWithJSON(w, http.StatusNoContent, nil)
}
func (h *databaseHandler) Clear(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.studentService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear students: %w", err))
		return
	}
	if err := h.teacherService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear teachers: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
