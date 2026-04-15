package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// LessonSlotHandler represents a handler responsible for LessonSlot-related HTTP operations.
type LessonSlotHandler interface {
	platform.BaseHandler[entities.LessonSlot]
}

// NewLessonSlotHandler creates a new LessonSlotHandler instance.
//
// It requires a lesson slot service.
func NewLessonSlotHandler(lss services.LessonSlotService) LessonSlotHandler {
	hc := platform.NewHandlerConfig("LessonSlotHandler", "lesson slot")

	return &lessonSlotHandler{
		BaseHandler: platform.NewBaseHandler(hc, lss),
		service:     lss,
	}
}

type lessonSlotHandler struct {
	platform.BaseHandler[entities.LessonSlot]
	service services.LessonSlotService
}
