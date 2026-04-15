package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// StudyLoadHandler represents a handler responsible for StudyLoad-related HTTP operations.
type StudyLoadHandler interface {
	platform.BaseHandler[entities.StudyLoad]
}

// NewStudyLoadHandler creates a new StudyLoadHandler instance.
//
// It requires a study load service.
func NewStudyLoadHandler(sls services.StudyLoadService) StudyLoadHandler {
	hc := platform.NewHandlerConfig("StudyLoadHandler", "study load")

	return &studyLoadHandler{
		BaseHandler: platform.NewBaseHandler(hc, sls),
		service:     sls,
	}
}

type studyLoadHandler struct {
	platform.BaseHandler[entities.StudyLoad]
	service services.StudyLoadService
}
