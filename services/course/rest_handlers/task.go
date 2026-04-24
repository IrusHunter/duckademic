package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/services/course/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type TaskHandler interface {
	platform.BaseHandler[entities.Task]
}

func NewTaskHandler(ts services.TaskService) TaskHandler {
	hc := platform.NewHandlerConfig("TaskHandler", entities.Task{}.EntityName())

	return &taskHandler{
		BaseHandler: platform.NewBaseHandler(hc, ts),
		service:     ts,
	}
}

type taskHandler struct {
	platform.BaseHandler[entities.Task]
	service services.TaskService
}
