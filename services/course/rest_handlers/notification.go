package resthandlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IrusHunter/duckademic/services/course/services"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/google/uuid"
)

type NotificationHandler interface {
	GetMyNotifications(context.Context, http.ResponseWriter, *http.Request)
	MarkAllRead(context.Context, http.ResponseWriter, *http.Request)
}

func NewNotificationHandler(ns services.NotificationService) NotificationHandler {
	return &notificationHandler{service: ns}
}

type notificationHandler struct {
	service services.NotificationService
}

func (h *notificationHandler) GetMyNotifications(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, err := h.userIDFromCtx(ctx)
	if err != nil {
		jsonutil.ResponseWithError(w, http.StatusUnauthorized, err)
		return
	}
	notifs, err := h.service.GetForUser(ctx, userID)
	if err != nil {
		jsonutil.ResponseWithError(w, http.StatusInternalServerError, err)
		return
	}
	jsonutil.ResponseWithJSON(w, http.StatusOK, notifs)
}

func (h *notificationHandler) MarkAllRead(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, err := h.userIDFromCtx(ctx)
	if err != nil {
		jsonutil.ResponseWithError(w, http.StatusUnauthorized, err)
		return
	}
	if err := h.service.MarkAllRead(ctx, userID); err != nil {
		jsonutil.ResponseWithError(w, http.StatusInternalServerError, err)
		return
	}
	jsonutil.ResponseWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *notificationHandler) userIDFromCtx(ctx context.Context) (uuid.UUID, error) {
	claims := contextutil.GetAccessClaims(ctx)
	if claims == nil {
		return uuid.Nil, fmt.Errorf("unauthorized")
	}
	id, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id: %w", err)
	}
	return id, nil
}
