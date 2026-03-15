package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/google/uuid"
)

type HandlerConfig struct {
	ClassName  string
	EntityName string
}

// NewHandlerConfig creates a new HandlerConfig instance.
//
// It requires a class name (cn) and an entity name (en).
func NewHandlerConfig(cn, en string) HandlerConfig {
	return HandlerConfig{
		ClassName:  cn,
		EntityName: en,
	}
}

// BaseHandler represents a handler responsible for entity-related HTTP operations.
type BaseHandler[T fmt.Stringer] interface {
	// GetAll returns a json with all entities.
	GetAll(context.Context, http.ResponseWriter, *http.Request)
	// Update handles HTTP request to update an entity by ID.
	Update(context.Context, http.ResponseWriter, *http.Request)
	// Delete handles HTTP request to delete an entity by ID.
	Delete(context.Context, http.ResponseWriter, *http.Request)
	// Add handles HTTP request to add a new entity.
	Add(context.Context, http.ResponseWriter, *http.Request)
	// Find handles HTTP request to find an entity. by ID
	Find(context.Context, http.ResponseWriter, *http.Request)

	// ParseID extracts and validates an entity ID from the request path.
	ParseID(ctx context.Context, w http.ResponseWriter, r *http.Request, method string) (uuid.UUID, bool)
	// DecodeEntity extracts and decodes an entity from the request body.
	DecodeEntity(ctx context.Context, w http.ResponseWriter, r *http.Request, method string) (T, bool)
}

// NewBaseHandler creates a new BaseHandler instance.
//
// It requires a academic rank services.
func NewBaseHandler[T fmt.Stringer](hc HandlerConfig, bs BaseService[T]) BaseHandler[T] {
	return &baseHandler[T]{
		HandlerConfig: hc,
		service:       bs,
		logger:        logger.NewLogger(hc.ClassName+".txt", hc.ClassName),
	}
}

type baseHandler[T fmt.Stringer] struct {
	HandlerConfig
	service BaseService[T]
	logger  logger.Logger
}

func (h *baseHandler[T]) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	res := h.service.GetAll(ctx)
	h.logger.Log(contextutil.GetTraceID(ctx), "GetAll",
		fmt.Sprintf("%d entities found", len(res)), logger.HandlerOperationSuccess)
	jsonutil.ResponseWithJSON(w, 200, res)
}
func (h *baseHandler[T]) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	entityID, ok := h.ParseID(ctx, w, r, "Update")
	if !ok {
		return
	}

	entity, ok := h.DecodeEntity(ctx, w, r, "Update")
	if !ok {
		return
	}

	updatedE, err := h.service.Update(ctx, entityID, entity)
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Update",
			fmt.Errorf("failed to update %s with id %q in service: %w", entity, entityID, err), logger.HandlerBadRequest),
		)
		return
	}

	h.logger.Log(contextutil.GetTraceID(ctx), "Update",
		fmt.Sprintf("%s successfully updated", updatedE), logger.HandlerOperationSuccess)
	jsonutil.ResponseWithJSON(w, 200, updatedE)
}
func (h *baseHandler[T]) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	entityID, ok := h.ParseID(ctx, w, r, "Delete")
	if !ok {
		return
	}

	entity, err := h.service.Delete(ctx, entityID)
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Delete",
			fmt.Errorf("failed to delete %s with id %q in service: %w", h.EntityName, entityID, err), logger.HandlerBadRequest),
		)
		return
	}

	h.logger.Log(contextutil.GetTraceID(ctx), "Delete",
		fmt.Sprintf("%s deleted", entity), logger.HandlerOperationSuccess)
	jsonutil.ResponseWithJSON(w, 200, entity)
}
func (h *baseHandler[T]) Add(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	entity, ok := h.DecodeEntity(ctx, w, r, "Add")
	if !ok {
		return
	}

	updatedE, err := h.service.Add(ctx, entity)
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Add",
			fmt.Errorf("failed to add %s in service: %w", entity, err), logger.HandlerBadRequest),
		)
		return
	}

	h.logger.Log(contextutil.GetTraceID(ctx), "Add",
		fmt.Sprintf("%s successfully added", updatedE), logger.HandlerOperationSuccess)
	jsonutil.ResponseWithJSON(w, 200, updatedE)
}
func (h *baseHandler[T]) Find(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	entityID, ok := h.ParseID(ctx, w, r, "Find")
	if !ok {
		return
	}

	entity := h.service.FindByID(ctx, entityID)
	if entity == nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "FindByID",
			fmt.Errorf("%s with id %q not found", h.EntityName, entityID), logger.HandlerBadRequest),
		)
		return
	}

	h.logger.Log(contextutil.GetTraceID(ctx), "FindByID",
		fmt.Sprintf("%s found", entity), logger.HandlerOperationSuccess)
	jsonutil.ResponseWithJSON(w, 200, entity)
}

func (h *baseHandler[T]) ParseID(ctx context.Context, w http.ResponseWriter, r *http.Request, method string) (uuid.UUID, bool) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx), method,
			fmt.Errorf("invalid id %q in the url path: %w", r.PathValue("id"), err),
			logger.HandlerBadRequest,
		))
		return uuid.Nil, false
	}
	return id, true
}
func (h *baseHandler[T]) DecodeEntity(ctx context.Context, w http.ResponseWriter, r *http.Request, method string) (T, bool) {
	var entity T
	err := json.NewDecoder(r.Body).Decode(&entity)
	defer r.Body.Close()

	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			method,
			fmt.Errorf("failed to extract %s from request body: %w", h.EntityName, err),
			logger.HandlerBadRequest,
		))
		return entity, false
	}

	return entity, true
}
