package resthandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/services/auth/services"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

// UserHandler represents a handler responsible for User-related HTTP operations.
type UserHandler interface {
	platform.BaseHandler[entities.User]
	ChangePassword(context.Context, http.ResponseWriter, *http.Request)
	Login(context.Context, http.ResponseWriter, *http.Request)
	Refresh(context.Context, http.ResponseWriter, *http.Request)
	ResetPassword(context.Context, http.ResponseWriter, *http.Request)
}

// NewUserHandler creates a new UserHandler instance.
//
// It requires a user service.
func NewUserHandler(us services.UserService) UserHandler {
	hc := platform.NewHandlerConfig("UserHandler", "user")

	return &userHandler{
		BaseHandler: platform.NewBaseHandler(hc, us),
		service:     us,
	}
}

type userHandler struct {
	platform.BaseHandler[entities.User]
	service services.UserService
}

func (h *userHandler) ChangePassword(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID          uuid.UUID `json:"id"`
		Login       string    `json:"login"`
		Password    string    `json:"password"`
		NewPassword string    `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonutil.ResponseWithError(w, 400, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"ChangePassword",
			fmt.Errorf("failed to decode new password: %w", err),
			logger.HandlerBadRequest,
		))
		return
	}
	r.Body.Close()

	err := h.service.ChangePassword(ctx, body.ID, body.Password, body.NewPassword)
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"ChangePassword",
			err,
			logger.HandlerBadRequest,
		))
		return
	}

	h.GetLogger().Log(contextutil.GetTraceID(ctx), "ChangePassword",
		fmt.Sprintf("password successfully changed for %s", body.Login),
		logger.HandlerOperationSuccess,
	)

	jsonutil.ResponseWithJSON(w, 204, nil)
}
func (h *userHandler) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var body struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonutil.ResponseWithError(w, 400, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"Login",
			fmt.Errorf("failed to decode credentials: %w", err),
			logger.HandlerBadRequest,
		))
		return
	}

	res, err := h.service.Login(ctx, body.Login, body.Password)
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"Login",
			err,
			logger.HandlerBadRequest,
		))
		return
	}

	h.GetLogger().Log(contextutil.GetTraceID(ctx), "Login",
		fmt.Sprintf("user %s logged in", res.Login),
		logger.HandlerOperationSuccess,
	)

	jsonutil.ResponseWithJSON(w, 200, res)
}
func (h *userHandler) Refresh(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var body struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonutil.ResponseWithError(w, 400, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"Refresh",
			fmt.Errorf("failed to decode tokens: %w", err),
			logger.HandlerBadRequest,
		))
		return
	}

	newAT, newRT, err := h.service.Refresh(ctx, body.AccessToken, body.RefreshToken)
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"Refresh",
			err,
			logger.HandlerBadRequest,
		))
		return
	}

	h.GetLogger().Log(contextutil.GetTraceID(ctx), "Refresh",
		"tokens successfully refreshed",
		logger.HandlerOperationSuccess,
	)

	jsonutil.ResponseWithJSON(w, 200, map[string]string{
		"access_token":  newAT,
		"refresh_token": newRT,
	})
}
func (h *userHandler) ResetPassword(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, ok := h.ParseID(ctx, w, r, "ResetPassword")
	if !ok {
		return
	}

	err := h.service.ResetPassword(ctx, userID)
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"ResetPassword",
			err,
			logger.HandlerBadRequest,
		))
		return
	}

	h.GetLogger().Log(contextutil.GetTraceID(ctx), "ResetPassword",
		fmt.Sprintf("password reset for user %s", userID),
		logger.HandlerOperationSuccess,
	)

	jsonutil.ResponseWithJSON(w, 204, nil)
}
