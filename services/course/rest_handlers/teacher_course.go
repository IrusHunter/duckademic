package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/services/course/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type TeacherCourseHandler interface {
	platform.BaseHandler[entities.TeacherCourse]
}

func NewTeacherCourseHandler(tcs services.TeacherCourseService) TeacherCourseHandler {
	hc := platform.NewHandlerConfig("TeacherCourseHandler", entities.TeacherCourse{}.EntityName())

	return &teacherCourseHandler{
		BaseHandler: platform.NewBaseHandler(hc, tcs),
		service:     tcs,
	}
}

type teacherCourseHandler struct {
	platform.BaseHandler[entities.TeacherCourse]
	service services.TeacherCourseService
}
