package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/curriculum/entities"
	"github.com/IrusHunter/duckademic/services/curriculum/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type CurriculumHandler interface {
	platform.BaseHandler[entities.Curriculum]
}

func NewCurriculumHandler(cs services.CurriculumService) CurriculumHandler {
	hc := platform.NewHandlerConfig("CurriculumHandler", "curriculum")

	return &curriculumHandler{
		BaseHandler: platform.NewBaseHandler(hc, cs),
		service:     cs,
	}
}

type curriculumHandler struct {
	platform.BaseHandler[entities.Curriculum]
	service services.CurriculumService
}
