package resthandlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/IrusHunter/duckademic/services/course/services"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
)

type DatabaseHandler interface {
	Seed(context.Context, http.ResponseWriter, *http.Request)
	Clear(context.Context, http.ResponseWriter, *http.Request)
}

func NewDatabaseHandler(
	ss services.StudentService,
	ts services.TeacherService,
	cs services.CourseService,
	scs services.StudentCourseService,
	tcs services.TeacherCourseService,
	taskS services.TaskService,
	tss services.TaskStudentService,
) DatabaseHandler {
	return &databaseHandler{
		studentService:       ss,
		teacherService:       ts,
		courseService:        cs,
		studentCourseService: scs,
		teacherCourseService: tcs,
		taskService:          taskS,
		taskStudentService:   tss,
	}
}

type databaseHandler struct {
	studentService       services.StudentService
	teacherService       services.TeacherService
	courseService        services.CourseService
	studentCourseService services.StudentCourseService
	teacherCourseService services.TeacherCourseService
	taskService          services.TaskService
	taskStudentService   services.TaskStudentService
}

func (h *databaseHandler) Seed(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(events.ExternalSeedCooldown * 2)
		ctx := contextutil.SetTraceID(context.Background())
		h.courseService.Seed(ctx)
		ctx = contextutil.SetTraceID(context.Background())
		h.teacherCourseService.Seed(ctx)
		ctx = contextutil.SetTraceID(context.Background())
		h.studentCourseService.Seed(ctx)
		ctx = contextutil.SetTraceID(context.Background())
		h.taskService.Seed(ctx)
		ctx = contextutil.SetTraceID(context.Background())
		h.taskStudentService.Seed(ctx)
	}()

	jsonutil.ResponseWithJSON(w, http.StatusNoContent, nil)
}
func (h *databaseHandler) Clear(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.taskStudentService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear task students: %w", err))
		return
	}
	if err := h.taskService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear tasks: %w", err))
		return
	}
	if err := h.studentCourseService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear student courses: %w", err))
		return
	}
	if err := h.teacherCourseService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear teacher courses: %w", err))
		return
	}
	if err := h.courseService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear courses: %w", err))
		return
	}
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
