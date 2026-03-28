package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type TeacherLoadHandler interface {
	platform.BaseHandler[entities.TeacherLoad]
}

func NewTeacherLoadHandler(ts services.TeacherLoadService) TeacherLoadHandler {
	hc := platform.NewHandlerConfig("TeacherLoadHandler", "teacher_load")

	return &teacherLoadHandler{
		BaseHandler: platform.NewBaseHandler(hc, ts),
		service:     ts,
	}
}

type teacherLoadHandler struct {
	platform.BaseHandler[entities.TeacherLoad]
	service services.TeacherLoadService
}
