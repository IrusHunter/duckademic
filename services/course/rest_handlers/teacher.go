package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/services/course/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// TeacherHandler represents a handler responsible for Teacher-related HTTP operations.
type TeacherHandler interface {
	platform.BaseHandler[entities.Teacher]
}

// NewTeacherHandler creates a new TeacherHandler instance.
//
// It requires a teacher service.
func NewTeacherHandler(ts services.TeacherService) TeacherHandler {
	hc := platform.NewHandlerConfig("TeacherHandler", entities.Teacher{}.EntityName())

	return &teacherHandler{
		BaseHandler: platform.NewBaseHandler(hc, ts),
		service:     ts,
	}
}

type teacherHandler struct {
	platform.BaseHandler[entities.Teacher]
	service services.TeacherService
}
