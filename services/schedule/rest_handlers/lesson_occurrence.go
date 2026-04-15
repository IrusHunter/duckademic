package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// LessonOccurrenceHandler represents a handler responsible for LessonOccurrence-related HTTP operations.
type LessonOccurrenceHandler interface {
	platform.BaseHandler[entities.LessonOccurrence]
}

// NewLessonOccurrenceHandler creates a new LessonOccurrenceHandler instance.
//
// It requires a lesson occurrence service.
func NewLessonOccurrenceHandler(los services.LessonOccurrenceService) LessonOccurrenceHandler {
	hc := platform.NewHandlerConfig("LessonOccurrenceHandler", entities.LessonOccurrence{}.EntityName())

	return &lessonOccurrenceHandler{
		BaseHandler: platform.NewBaseHandler(hc, los),
		service:     los,
	}
}

type lessonOccurrenceHandler struct {
	platform.BaseHandler[entities.LessonOccurrence]
	service services.LessonOccurrenceService
}
