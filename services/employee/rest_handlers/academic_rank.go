package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/employee/entities"
	"github.com/IrusHunter/duckademic/services/employee/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// AcademicRankHandler represents a handler responsible for AcademicRank-related HTTP operations.
type AcademicRankHandler interface {
	platform.BaseHandler[entities.AcademicRank]
}

// NewAcademicRankHandler creates a new AcademicRankHandler instance.
//
// It requires an academic rank services.
func NewAcademicRankHandler(ars services.AcademicRankService) AcademicRankHandler {
	hc := platform.NewHandlerConfig("AcademicRankHandler", "academic rank")

	return &academicRankHandler{
		BaseHandler: platform.NewBaseHandler(hc, ars),
		service:     ars,
	}
}

type academicRankHandler struct {
	platform.BaseHandler[entities.AcademicRank]
	service services.AcademicRankService
}
