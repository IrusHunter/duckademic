package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/services/course/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type TaskStudentHandler interface {
	platform.BaseHandler[entities.TaskStudent]
}

func NewTaskStudentHandler(tss services.TaskStudentService) TaskStudentHandler {
	hc := platform.NewHandlerConfig("TaskStudentHandler", entities.TaskStudent{}.EntityName())

	return &taskStudentHandler{
		BaseHandler: platform.NewBaseHandler(hc, tss),
		service:     tss,
	}
}

type taskStudentHandler struct {
	platform.BaseHandler[entities.TaskStudent]
	service services.TaskStudentService
}
