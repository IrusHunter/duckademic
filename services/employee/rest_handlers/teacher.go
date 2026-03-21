package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/employee/entities"
	"github.com/IrusHunter/duckademic/services/employee/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type TeacherHandler interface {
	platform.BaseHandler[entities.Teacher]
}

func NewTeacherHandler(ts services.TeacherService) TeacherHandler {
	hc := platform.NewHandlerConfig("TeacherHandler", "teacher")

	return &teacherHandler{
		BaseHandler: platform.NewBaseHandler(hc, ts),
		service:     ts,
	}
}

type teacherHandler struct {
	platform.BaseHandler[entities.Teacher]
	service services.TeacherService
}
