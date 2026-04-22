package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type SemesterHandler interface {
	platform.BaseHandler[entities.Semester]
}

func NewSemesterHandler(ss services.SemesterService) SemesterHandler {
	hc := platform.NewHandlerConfig("SemesterHandler", "semester")

	return &semesterHandler{
		BaseHandler: platform.NewBaseHandler(hc, ss),
		service:     ss,
	}
}

type semesterHandler struct {
	platform.BaseHandler[entities.Semester]
	service services.SemesterService
}
