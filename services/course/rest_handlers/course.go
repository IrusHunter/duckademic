package resthandlers

import (
	"context"
	"net/http"

	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/services/course/services"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type CourseHandler interface {
	platform.BaseHandler[entities.Course]
	GetStudentCoursePage(ctx context.Context, w http.ResponseWriter, r *http.Request)
	GetTeacherCoursePage(ctx context.Context, w http.ResponseWriter, r *http.Request)
}

func NewCourseHandler(cs services.CourseService) CourseHandler {
	hc := platform.NewHandlerConfig("CourseHandler", entities.Course{}.EntityName())

	return &courseHandler{
		BaseHandler: platform.NewBaseHandler(hc, cs),
		service:     cs,
	}
}

type courseHandler struct {
	platform.BaseHandler[entities.Course]
	service services.CourseService
}

func (h *courseHandler) GetStudentCoursePage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	studentID, ok := h.GetUserIDFromContext(ctx, w, "GetStudentCoursePage")
	if !ok {
		return
	}

	coursePages, err := h.service.GetStudentCoursePage(ctx, studentID)
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetStudentCoursePage",
			err,
			logger.HandlerInternalError,
		))
		return
	}

	h.GetLogger().Log(
		contextutil.GetTraceID(ctx),
		"GetStudentCoursePage",
		"student course pages fetched successfully",
		logger.HandlerOperationSuccess,
	)

	jsonutil.ResponseWithJSON(w, 200, coursePages)
}

func (h *courseHandler) GetTeacherCoursePage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	teacherID, ok := h.GetUserIDFromContext(ctx, w, "GetTeacherCoursePage")
	if !ok {
		return
	}

	coursePages, err := h.service.GetTeacherCoursePage(ctx, teacherID)
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetTeacherCoursePage",
			err,
			logger.HandlerInternalError,
		))
		return
	}

	h.GetLogger().Log(
		contextutil.GetTraceID(ctx),
		"GetTeacherCoursePage",
		"teacher course pages fetched successfully",
		logger.HandlerOperationSuccess,
	)

	jsonutil.ResponseWithJSON(w, 200, coursePages)
}
