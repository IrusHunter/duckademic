package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/services"
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
