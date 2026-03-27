package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/teacher_load/entities"
	"github.com/IrusHunter/duckademic/services/teacher_load/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type DisciplineHandler interface {
	platform.BaseHandler[entities.Discipline]
}

func NewDisciplineHandler(ds services.DisciplineService) DisciplineHandler {
	hc := platform.NewHandlerConfig("DisciplineHandler", "discipline")

	return &disciplineHandler{
		BaseHandler: platform.NewBaseHandler(hc, ds),
		service:     ds,
	}
}

type disciplineHandler struct {
	platform.BaseHandler[entities.Discipline]
	service services.DisciplineService
}
