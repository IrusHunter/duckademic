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
	lts services.LessonTypeService,
	ds services.DisciplineService,
	ltas services.LessonTypeAssignmentService,
) DatabaseHandler {
	return &databaseHandler{
		curriculumService:    cs,
		semesterService:      ss,
		lessonTypeService:    lts,
		disciplineService:    ds,
		lessonTypeAssignment: ltas,
	}
}

type databaseHandler struct {
	curriculumService    services.CurriculumService
	semesterService      services.SemesterService
	lessonTypeService    services.LessonTypeService
	disciplineService    services.DisciplineService
	lessonTypeAssignment services.LessonTypeAssignmentService
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
	if err := h.lessonTypeService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed lesson types: %w", err))
		return
	}
	if err := h.disciplineService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed disciplines: %w", err))
		return
	}
	if err := h.lessonTypeAssignment.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed lesson type assignments: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
func (h *databaseHandler) Clear(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.lessonTypeAssignment.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear lesson type assignments: %w", err))
		return
	}
	if err := h.semesterService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear semesters: %w", err))
		return
	}
	if err := h.curriculumService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear curriculums: %w", err))
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

	jsonutil.ResponseWithJSON(w, 204, nil)
}
