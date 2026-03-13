package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/employees/entities"
	"github.com/IrusHunter/duckademic/services/employees/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// AcademicDegreeHandler represents a handler responsible for AcademicDegree-related HTTP operations.
type AcademicDegreeHandler interface {
	platform.BaseHandler[entities.AcademicDegree]
}

// NewAcademicRankHandler creates a new AcademicDegreeHandler instance.
//
// It requires an academic degree services.
func NewAcademicDegreeHandler(ads services.AcademicDegreeService) AcademicDegreeHandler {
	hc := platform.NewHandlerConfig("AcademicRankHandler", "academic rank")

	return &academicDegreeHandler{
		BaseHandler: platform.NewBaseHandler(hc, ads),
		service:     ads,
	}
}

type academicDegreeHandler struct {
	platform.BaseHandler[entities.AcademicDegree]
	service services.AcademicDegreeService
}
