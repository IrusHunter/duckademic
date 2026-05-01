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

type TaskStudentHandler interface {
	platform.BaseHandler[entities.TaskStudent]
	GetUpcomingEvents(context.Context, http.ResponseWriter, *http.Request)
}

func NewTaskStudentHandler(tss services.TaskStudentService) TaskStudentHandler {
	hc := platform.NewHandlerConfig("TaskStudentHandler", entities.TaskStudent{}.EntityName())

	return &taskStudentHandler{
		BaseHandler: platform.NewBaseHandler(hc, tss),
		service:     tss,
	}
}

type taskStudentHandler struct {
	platform.BaseHandler[entities.TaskStudent]
	service services.TaskStudentService
}

func (h *taskStudentHandler) GetUpcomingEvents(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	count, ok := h.ParseIntQueryParam(ctx, w, q, "count", "GetUpcomingEventsFor")
	if !ok {
		return
	}
	startTime, ok := h.ParseTimeQueryParam(ctx, w, q, "start-time", "GetUpcomingEventsFor")
	if !ok {
		return
	}
	studentID, ok := h.GetUserIDFromContext(ctx, w, "GetUpcomingTasksFor")
	if !ok {
		return
	}

	tasks, err := h.service.GetUpcomingTasksFor(ctx, studentID, startTime, count)
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetUpcomingTasksFor",
			err,
			logger.HandlerInternalError,
		))
		return
	}

	h.GetLogger().Log(
		contextutil.GetTraceID(ctx),
		"GetUpcomingTasksFor",
		"upcoming tasks fetched successfully",
		logger.HandlerOperationSuccess,
	)

	jsonutil.ResponseWithJSON(w, 200, tasks)
}
