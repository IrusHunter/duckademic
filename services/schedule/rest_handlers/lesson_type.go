package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type LessonTypeHandler interface {
	platform.BaseHandler[entities.LessonType]
}

func NewLessonTypeHandler(ls services.LessonTypeService) LessonTypeHandler {
	hc := platform.NewHandlerConfig("LessonTypeHandler", "lesson_type")

	return &lessonTypeHandler{
		BaseHandler: platform.NewBaseHandler(hc, ls),
		service:     ls,
	}
}

type lessonTypeHandler struct {
	platform.BaseHandler[entities.LessonType]
	service services.LessonTypeService
}
