package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// TeacherUnavailableDayHandler represents HTTP operations for TeacherUnavailableDays.
type TeacherUnavailableDayHandler interface {
	platform.BaseHandler[entities.TeacherUnavailableDay]
}

// NewTeacherUnavailableDayHandler creates a new handler instance.
func NewTeacherUnavailableDayHandler(
	s services.TeacherUnavailableDayService,
) TeacherUnavailableDayHandler {
	hc := platform.NewHandlerConfig(
		"TeacherUnavailableDayHandler",
		entities.TeacherUnavailableDay{}.EntityName(),
	)

	return &teacherUnavailableDayHandler{
		BaseHandler: platform.NewBaseHandler(hc, s),
		service:     s,
	}
}

type teacherUnavailableDayHandler struct {
	platform.BaseHandler[entities.TeacherUnavailableDay]
	service services.TeacherUnavailableDayService
}
