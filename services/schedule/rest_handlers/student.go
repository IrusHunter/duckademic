package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type StudentHandler interface {
	platform.BaseHandler[entities.Student]
}

func NewStudentHandler(ss services.StudentService) StudentHandler {
	hc := platform.NewHandlerConfig("StudentHandler", "student")

	return &studentHandler{
		BaseHandler: platform.NewBaseHandler(hc, ss),
		service:     ss,
	}
}

type studentHandler struct {
	platform.BaseHandler[entities.Student]
	service services.StudentService
}
