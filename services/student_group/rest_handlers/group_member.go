package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/student_group/entities"
	"github.com/IrusHunter/duckademic/services/student_group/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type GroupMembersHandler interface {
	platform.BaseHandler[entities.GroupMember]
}

func NewGroupMembersHandler(gms services.GroupMemberService) GroupMembersHandler {
	hc := platform.NewHandlerConfig("GroupMembersHandler", "group member")

	return &groupMembersHandler{
		BaseHandler: platform.NewBaseHandler(hc, gms),
		service:     gms,
	}
}

type groupMembersHandler struct {
	platform.BaseHandler[entities.GroupMember]
	service services.GroupMemberService
}
