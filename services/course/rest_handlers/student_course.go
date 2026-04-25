package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/services/course/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// StudentCourseHandler represents a handler responsible for StudentCourse-related HTTP operations.
type StudentCourseHandler interface {
	platform.BaseHandler[entities.StudentCourse]
}

// NewStudentCourseHandler creates a new StudentCourseHandler instance.
func NewStudentCourseHandler(scs services.StudentCourseService) StudentCourseHandler {
	hc := platform.NewHandlerConfig("StudentCourseHandler", entities.StudentCourse{}.EntityName())

	return &studentCourseHandler{
		BaseHandler: platform.NewBaseHandler(hc, scs),
		service:     scs,
	}
}

type studentCourseHandler struct {
	platform.BaseHandler[entities.StudentCourse]
	service services.StudentCourseService
}
