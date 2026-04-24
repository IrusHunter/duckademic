package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/services/course/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type CourseHandler interface {
	platform.BaseHandler[entities.Course]
}

func NewCourseHandler(cs services.CourseService) CourseHandler {
	hc := platform.NewHandlerConfig("CourseHandler", entities.Course{}.EntityName())

	return &courseHandler{
		BaseHandler: platform.NewBaseHandler(hc, cs),
		service:     cs,
	}
}

type courseHandler struct {
	platform.BaseHandler[entities.Course]
	service services.CourseService
}
