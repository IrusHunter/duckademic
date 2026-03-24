package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/curriculum/entities"
	"github.com/IrusHunter/duckademic/services/curriculum/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type DisciplineHandler interface {
	platform.BaseHandler[entities.Discipline]
}

func NewDisciplineHandler(ds services.DisciplineService) DisciplineHandler {
	hd := platform.NewHandlerConfig("DisciplineHandler", "discipline")

	return &disciplineHandler{
		BaseHandler: platform.NewBaseHandler(hd, ds),
		service:     ds,
	}
}

type disciplineHandler struct {
	platform.BaseHandler[entities.Discipline]
	service services.DisciplineService
}
