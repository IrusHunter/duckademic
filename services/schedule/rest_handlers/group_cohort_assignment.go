package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type GroupCohortAssignmentHandler interface {
	platform.BaseHandler[entities.GroupCohortAssignment]
}

func NewGroupCohortAssignmentHandler(gcas services.GroupCohortAssignmentService) GroupCohortAssignmentHandler {
	hc := platform.NewHandlerConfig("GroupCohortAssignmentHandler", "group cohort assignment")

	return &groupCohortAssignmentHandler{
		BaseHandler: platform.NewBaseHandler(hc, gcas),
		service:     gcas,
	}
}

type groupCohortAssignmentHandler struct {
	platform.BaseHandler[entities.GroupCohortAssignment]
	service services.GroupCohortAssignmentService
}
