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

type TaskStudentHandler interface {
	platform.BaseHandler[entities.TaskStudent]
	GetUpcomingEvents(context.Context, http.ResponseWriter, *http.Request)
	GetForStudentInCourse(context.Context, http.ResponseWriter, *http.Request)
	GetByTaskID(context.Context, http.ResponseWriter, *http.Request)
}

func NewTaskStudentHandler(tss services.TaskStudentService, ns services.NotificationService) TaskStudentHandler {
	hc := platform.NewHandlerConfig("TaskStudentHandler", entities.TaskStudent{}.EntityName())

	return &taskStudentHandler{
		BaseHandler:  platform.NewBaseHandler(hc, tss),
		service:      tss,
		notifService: ns,
	}
}

type taskStudentHandler struct {
	platform.BaseHandler[entities.TaskStudent]
	service      services.TaskStudentService
	notifService services.NotificationService
}

func (h *taskStudentHandler) Add(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

	go h.notifService.NotifySubmission(contextutil.SetTraceID(context.Background()), created.TaskID, created.StudentID)

	h.GetLogger().Log(contextutil.GetTraceID(ctx), "Add",
		fmt.Sprintf("%s successfully added", created), logger.HandlerOperationSuccess)
	jsonutil.ResponseWithJSON(w, http.StatusOK, created)
}

func (h *taskStudentHandler) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	entityID, ok := h.ParseID(ctx, w, r, "Update")
	if !ok {
		return
	}

	existing := h.service.FindByID(ctx, entityID)

	entity, ok := h.DecodeEntity(ctx, w, r, "Update")
	if !ok {
		return
	}

	updated, err := h.service.Update(ctx, entityID, entity)
	if err != nil {
		jsonutil.ResponseWithError(w, http.StatusBadRequest, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx), "Update", err, logger.HandlerBadRequest,
		))
		return
	}

	if existing != nil && existing.Mark == nil && entity.Mark != nil {
		go h.notifService.NotifyGrade(contextutil.SetTraceID(context.Background()), existing.TaskID, existing.StudentID)
	}

	h.GetLogger().Log(contextutil.GetTraceID(ctx), "Update",
		fmt.Sprintf("%s successfully updated", updated), logger.HandlerOperationSuccess)
	jsonutil.ResponseWithJSON(w, http.StatusOK, updated)
}

func (h *taskStudentHandler) GetForStudentInCourse(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	studentID, ok := h.GetUserIDFromContext(ctx, w, "GetForStudentInCourse")
	if !ok {
		return
	}

	courseID, ok := h.ParseID(ctx, w, r, "GetForStudentInCourse")
	if !ok {
		return
	}

	submissions, err := h.service.GetForStudentInCourse(ctx, studentID, courseID)
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetForStudentInCourse",
			err,
			logger.HandlerInternalError,
		))
		return
	}

	h.GetLogger().Log(
		contextutil.GetTraceID(ctx),
		"GetForStudentInCourse",
		fmt.Sprintf("%d submissions found", len(submissions)),
		logger.HandlerOperationSuccess,
	)

	jsonutil.ResponseWithJSON(w, 200, submissions)
}

func (h *taskStudentHandler) GetByTaskID(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	taskID, ok := h.ParseID(ctx, w, r, "GetByTaskID")
	if !ok {
		return
	}

	submissions, err := h.service.GetByTaskID(ctx, taskID)
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx), "GetByTaskID", err, logger.HandlerInternalError,
		))
		return
	}

	h.GetLogger().Log(contextutil.GetTraceID(ctx), "GetByTaskID",
		fmt.Sprintf("%d submissions found for task", len(submissions)),
		logger.HandlerOperationSuccess,
	)

	jsonutil.ResponseWithJSON(w, 200, submissions)
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
