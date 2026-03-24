package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/curriculum/entities"
	"github.com/IrusHunter/duckademic/services/curriculum/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type SemesterDisciplineHandler interface {
	platform.BaseHandler[entities.SemesterDiscipline]
}

func NewSemesterDisciplineHandler(sds services.SemesterDisciplineService) SemesterDisciplineHandler {
	hc := platform.NewHandlerConfig("SemesterDisciplineHandler", "semester_discipline")

	return &semesterDisciplineHandler{
		BaseHandler: platform.NewBaseHandler(hc, sds),
		service:     sds,
	}
}

type semesterDisciplineHandler struct {
	platform.BaseHandler[entities.SemesterDiscipline]
	service services.SemesterDisciplineService
}
