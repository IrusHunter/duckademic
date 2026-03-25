package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type LessonTypeAssignmentHandler interface {
	platform.BaseHandler[entities.LessonTypeAssignment]
}

func NewLessonTypeAssignmentHandler(ltaService services.LessonTypeAssignmentService) LessonTypeAssignmentHandler {
	hc := platform.NewHandlerConfig("LessonTypeAssignmentHandler", "lesson_type_assignment")

	return &lessonTypeAssignmentHandler{
		BaseHandler: platform.NewBaseHandler(hc, ltaService),
		service:     ltaService,
	}
}

type lessonTypeAssignmentHandler struct {
	platform.BaseHandler[entities.LessonTypeAssignment]
	service services.LessonTypeAssignmentService
}
