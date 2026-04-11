package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/asset/entities"
	"github.com/IrusHunter/duckademic/services/asset/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type ClassroomHandler interface {
	platform.BaseHandler[entities.Classroom]
}

func NewClassroomHandler(cs services.ClassroomService) ClassroomHandler {
	hc := platform.NewHandlerConfig("ClassroomHandler", "classroom")

	return &classroomHandler{
		BaseHandler: platform.NewBaseHandler(hc, cs),
		service:     cs,
	}
}

type classroomHandler struct {
	platform.BaseHandler[entities.Classroom]
	service services.ClassroomService
}
