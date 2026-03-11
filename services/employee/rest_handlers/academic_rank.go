package resthandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IrusHunter/duckademic/services/employees/entities"
	"github.com/IrusHunter/duckademic/services/employees/services"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/google/uuid"
)

// AcademicRankHandler represents a handler responsible for AcademicRank-related HTTP operations.
type AcademicRankHandler interface {
	// GetAll returns a json with all academic ranks.
	GetAll(context.Context, http.ResponseWriter, *http.Request)
	// Update handles HTTP request to update an AcademicRank by ID.
	Update(context.Context, http.ResponseWriter, *http.Request)
	// Delete handles HTTP request to delete an AcademicRank by ID.
	Delete(context.Context, http.ResponseWriter, *http.Request)
	// Add handles HTTP request to add a new AcademicRank.
	Add(context.Context, http.ResponseWriter, *http.Request)
	// Find handles HTTP request to find AcademicRank.
	Find(context.Context, http.ResponseWriter, *http.Request)
}

// NewAcademicRankHandler creates a new AcademicRankHandler instance.
//
// It requires a academic rank services.
func NewAcademicRankHandler(ars services.AcademicRankService) AcademicRankHandler {
	return &academicRankHandler{service: ars}
}

type academicRankHandler struct {
	service services.AcademicRankService
}

func (h *academicRankHandler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	jsonutil.ResponseWithJSON(w, 200, h.service.GetAll(ctx))
}
func (h *academicRankHandler) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	academicRankID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		jsonutil.ResponseWithError(w, 400, fmt.Errorf("invalid id %q in url path: %w", r.PathValue("id"), err))
		return
	}

	academicRank := entities.AcademicRank{}
	err = json.NewDecoder(r.Body).Decode(&academicRank)
	defer r.Body.Close()

	if err != nil {
		jsonutil.ResponseWithError(w, 400, fmt.Errorf("failed to extract academic rank from request body: %w", err))
		return
	}

	academicRank, err = h.service.Update(ctx, academicRankID, academicRank)
	if err != nil {
		jsonutil.ResponseWithError(w, 400, fmt.Errorf("failed to update in service: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 200, academicRank)
}
func (h *academicRankHandler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	academicRankID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		jsonutil.ResponseWithError(w, 400, fmt.Errorf("invalid id %q in url path: %w", r.PathValue("id"), err))
		return
	}

	err = h.service.Delete(ctx, academicRankID)
	if err != nil {
		jsonutil.ResponseWithError(w, 400, fmt.Errorf("failed to delete in service: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
func (h *academicRankHandler) Add(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	academicRank := entities.AcademicRank{}
	err := json.NewDecoder(r.Body).Decode(&academicRank)
	defer r.Body.Close()

	if err != nil {
		jsonutil.ResponseWithError(w, 400, fmt.Errorf("failed to extract academic rank from request body: %w", err))
		return
	}

	academicRank, err = h.service.Add(ctx, academicRank)
	if err != nil {
		jsonutil.ResponseWithError(w, 400, fmt.Errorf("failed to add in service: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 200, academicRank)
}
func (h *academicRankHandler) Find(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	academicRankID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		jsonutil.ResponseWithError(w, 400, fmt.Errorf("invalid id %q in url path: %w", r.PathValue("id"), err))
		return
	}

	academicRank := h.service.FindByID(ctx, academicRankID)
	if academicRank == nil {
		jsonutil.ResponseWithError(w, 400, fmt.Errorf("academic rank with id %q not found", academicRankID))
		return
	}

	jsonutil.ResponseWithJSON(w, 200, academicRank)
}
