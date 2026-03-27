package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/teacher_load/entities"
	"github.com/IrusHunter/duckademic/services/teacher_load/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type GroupCohortHandler interface {
	platform.BaseHandler[entities.GroupCohort]
}

func NewGroupCohortHandler(gc services.GroupCohortService) GroupCohortHandler {
	hc := platform.NewHandlerConfig("GroupCohortHandler", "group cohort")

	return &groupCohortHandler{
		BaseHandler: platform.NewBaseHandler(hc, gc),
		service:     gc,
	}
}

type groupCohortHandler struct {
	platform.BaseHandler[entities.GroupCohort]
	service services.GroupCohortService
}
