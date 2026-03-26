package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type StudentGroupHandler interface {
	platform.BaseHandler[entities.StudentGroup]
}

func NewStudentGroupHandler(sg services.StudentGroupService) StudentGroupHandler {
	hc := platform.NewHandlerConfig("StudentGroupHandler", "student group")

	return &studentGroupHandler{
		BaseHandler: platform.NewBaseHandler(hc, sg),
		service:     sg,
	}
}

type studentGroupHandler struct {
	platform.BaseHandler[entities.StudentGroup]
	service services.StudentGroupService
}
