package resthandlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/services/course/services"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type TaskHandler interface {
	platform.BaseHandler[entities.Task]
	GetByCourseID(context.Context, http.ResponseWriter, *http.Request)
}

func NewTaskHandler(ts services.TaskService, ns services.NotificationService) TaskHandler {
	hc := platform.NewHandlerConfig("TaskHandler", entities.Task{}.EntityName())

	return &taskHandler{
		BaseHandler:  platform.NewBaseHandler(hc, ts),
		service:      ts,
		notifService: ns,
	}
}

type taskHandler struct {
	platform.BaseHandler[entities.Task]
	service      services.TaskService
	notifService services.NotificationService
}

func (h *taskHandler) Add(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	entity, ok := h.DecodeEntity(ctx, w, r, "Add")
	if !ok {
		return
	}

	created, err := h.service.Add(ctx, entity)
	if err != nil {
		jsonutil.ResponseWithError(w, http.StatusBadRequest, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx), "Add", err, logger.HandlerBadRequest,
		))
		return
	}

	go h.notifService.NotifyNewTask(contextutil.SetTraceID(context.Background()), created)

	h.GetLogger().Log(contextutil.GetTraceID(ctx), "Add",
		fmt.Sprintf("%s successfully added", created), logger.HandlerOperationSuccess)
	jsonutil.ResponseWithJSON(w, http.StatusOK, created)
}

func (h *taskHandler) GetByCourseID(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	courseID, ok := h.ParseID(ctx, w, r, "GetByCourseID")
	if !ok {
		return
	}

	tasks, err := h.service.GetByCourseID(ctx, courseID)
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetByCourseID",
			err,
			logger.HandlerInternalError,
		))
		return
	}

	h.GetLogger().Log(
		contextutil.GetTraceID(ctx),
		"GetByCourseID",
		fmt.Sprintf("%d tasks found for course %s", len(tasks), courseID),
		logger.HandlerOperationSuccess,
	)

	jsonutil.ResponseWithJSON(w, 200, tasks)
}
