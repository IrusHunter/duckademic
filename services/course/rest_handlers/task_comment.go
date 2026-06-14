package resthandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/services/course/services"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type TaskCommentHandler interface {
	platform.BaseHandler[entities.TaskComment]
	GetForTask(context.Context, http.ResponseWriter, *http.Request)
	GetPrivateComments(context.Context, http.ResponseWriter, *http.Request)
	AddComment(context.Context, http.ResponseWriter, *http.Request)
	AddPrivateComment(context.Context, http.ResponseWriter, *http.Request)
}

func NewTaskCommentHandler(svc services.TaskCommentService) TaskCommentHandler {
	hc := platform.NewHandlerConfig("TaskCommentHandler", entities.TaskComment{}.EntityName())
	return &taskCommentHandler{
		BaseHandler: platform.NewBaseHandler(hc, svc),
		service:     svc,
	}
}

type taskCommentHandler struct {
	platform.BaseHandler[entities.TaskComment]
	service services.TaskCommentService
}

type addCommentRequest struct {
	Body      string     `json:"body"`
	StudentID *uuid.UUID `json:"student_id,omitempty"`
}

func (h *taskCommentHandler) GetForTask(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	taskID, ok := h.ParseID(ctx, w, r, "GetForTask")
	if !ok {
		return
	}

	comments, err := h.service.GetForTask(ctx, taskID, false, nil)
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx), "GetForTask", err, logger.HandlerInternalError,
		))
		return
	}

	h.GetLogger().Log(contextutil.GetTraceID(ctx), "GetForTask",
		fmt.Sprintf("%d class comments found", len(comments)), logger.HandlerOperationSuccess)
	jsonutil.ResponseWithJSON(w, 200, comments)
}

func (h *taskCommentHandler) GetPrivateComments(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	taskID, ok := h.ParseID(ctx, w, r, "GetPrivateComments")
	if !ok {
		return
	}

	claims := contextutil.GetAccessClaims(ctx)
	if claims == nil {
		jsonutil.ResponseWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	var studentID uuid.UUID
	if claims.Role == "teacher" {
		sidStr := r.URL.Query().Get("student_id")
		if sidStr == "" {
			jsonutil.ResponseWithError(w, 400, fmt.Errorf("student_id query param required for teacher"))
			return
		}
		parsed, err := uuid.Parse(sidStr)
		if err != nil {
			jsonutil.ResponseWithError(w, 400, fmt.Errorf("invalid student_id"))
			return
		}
		studentID = parsed
	} else {
		uid, err := uuid.Parse(claims.UserID)
		if err != nil {
			jsonutil.ResponseWithError(w, http.StatusUnauthorized, fmt.Errorf("invalid user id in token"))
			return
		}
		studentID = uid
	}

	comments, err := h.service.GetForTask(ctx, taskID, true, &studentID)
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx), "GetPrivateComments", err, logger.HandlerInternalError,
		))
		return
	}

	h.GetLogger().Log(contextutil.GetTraceID(ctx), "GetPrivateComments",
		fmt.Sprintf("%d private comments found", len(comments)), logger.HandlerOperationSuccess)
	jsonutil.ResponseWithJSON(w, 200, comments)
}

func (h *taskCommentHandler) AddComment(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	taskID, ok := h.ParseID(ctx, w, r, "AddComment")
	if !ok {
		return
	}

	authorID, ok := h.GetUserIDFromContext(ctx, w, "AddComment")
	if !ok {
		return
	}

	claims := contextutil.GetAccessClaims(ctx)
	if claims == nil {
		jsonutil.ResponseWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	var req addCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Body == "" {
		jsonutil.ResponseWithError(w, 400, fmt.Errorf("body field is required"))
		return
	}

	comment, err := h.service.AddComment(ctx, taskID, authorID, claims.Role, req.Body, false, nil)
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx), "AddComment", err, logger.HandlerBadRequest,
		))
		return
	}

	h.GetLogger().Log(contextutil.GetTraceID(ctx), "AddComment",
		"class comment added", logger.HandlerOperationSuccess)
	jsonutil.ResponseWithJSON(w, 200, comment)
}

func (h *taskCommentHandler) AddPrivateComment(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	taskID, ok := h.ParseID(ctx, w, r, "AddPrivateComment")
	if !ok {
		return
	}

	authorID, ok := h.GetUserIDFromContext(ctx, w, "AddPrivateComment")
	if !ok {
		return
	}

	claims := contextutil.GetAccessClaims(ctx)
	if claims == nil {
		jsonutil.ResponseWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	var req addCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Body == "" {
		jsonutil.ResponseWithError(w, 400, fmt.Errorf("body field is required"))
		return
	}

	var studentID *uuid.UUID
	if claims.Role == "teacher" {
		if req.StudentID == nil {
			jsonutil.ResponseWithError(w, 400, fmt.Errorf("student_id required for teacher private comment"))
			return
		}
		studentID = req.StudentID
	} else {
		studentID = &authorID
	}

	comment, err := h.service.AddComment(ctx, taskID, authorID, claims.Role, req.Body, true, studentID)
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx), "AddPrivateComment", err, logger.HandlerBadRequest,
		))
		return
	}

	h.GetLogger().Log(contextutil.GetTraceID(ctx), "AddPrivateComment",
		"private comment added", logger.HandlerOperationSuccess)
	jsonutil.ResponseWithJSON(w, 200, comment)
}
