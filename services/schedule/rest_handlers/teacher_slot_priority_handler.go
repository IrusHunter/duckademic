package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// TeacherSlotPriorityHandler represents HTTP operations for TeacherSlotPriority.
type TeacherSlotPriorityHandler interface {
	platform.BaseHandler[entities.TeacherSlotPriority]
}

// NewTeacherSlotPriorityHandler creates a new handler instance.
func NewTeacherSlotPriorityHandler(
	s services.TeacherSlotPriorityService,
) TeacherSlotPriorityHandler {
	hc := platform.NewHandlerConfig(
		"TeacherSlotPriorityHandler",
		entities.TeacherSlotPriority{}.EntityName(),
	)

	return &teacherSlotPriorityHandler{
		BaseHandler: platform.NewBaseHandler(hc, s),
		service:     s,
	}
}

type teacherSlotPriorityHandler struct {
	platform.BaseHandler[entities.TeacherSlotPriority]
	service services.TeacherSlotPriorityService
}
