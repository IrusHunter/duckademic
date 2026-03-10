package resthandlers

import (
	"context"
	"net/http"

	"github.com/IrusHunter/duckademic/services/employees/services"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
)

// AcademicRankHandler represents a handler responsible for AcademicRank-related HTTP operations.
type AcademicRankHandler interface {
	// GetAll returns a json with all academic ranks.
	GetAll(context.Context, http.ResponseWriter, *http.Request)
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
