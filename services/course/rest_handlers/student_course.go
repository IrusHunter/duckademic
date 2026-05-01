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

// StudentCourseHandler represents a handler responsible for StudentCourse-related HTTP operations.
type StudentCourseHandler interface {
	platform.BaseHandler[entities.StudentCourse]
	GetCourseProgress(context.Context, http.ResponseWriter, *http.Request)
}

// NewStudentCourseHandler creates a new StudentCourseHandler instance.
func NewStudentCourseHandler(scs services.StudentCourseService) StudentCourseHandler {
	hc := platform.NewHandlerConfig("StudentCourseHandler", entities.StudentCourse{}.EntityName())

	return &studentCourseHandler{
		BaseHandler: platform.NewBaseHandler(hc, scs),
		service:     scs,
	}
}

type studentCourseHandler struct {
	platform.BaseHandler[entities.StudentCourse]
	service services.StudentCourseService
}

func (h *studentCourseHandler) GetCourseProgress(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	studentID, ok := h.GetUserIDFromContext(ctx, w, "GetCourseProgress")
	if !ok {
		return
	}

	courseProgress, err := h.service.GetCourseProgress(ctx, studentID)
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetCourseProgress",
			err,
			logger.HandlerInternalError,
		))
		return
	}

	h.GetLogger().Log(
		contextutil.GetTraceID(ctx),
		"GetCourseProgress",
		"courses progress fetched successfully",
		logger.HandlerOperationSuccess,
	)

	jsonutil.ResponseWithJSON(w, 200, courseProgress)
}
