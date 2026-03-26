package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type GroupMemberHandler interface {
	platform.BaseHandler[entities.GroupMember]
}

func NewGroupMemberHandler(gms services.GroupMemberService) GroupMemberHandler {
	hc := platform.NewHandlerConfig("GroupMembersHandler", "group member")

	return &groupMemberHandler{
		BaseHandler: platform.NewBaseHandler(hc, gms),
		service:     gms,
	}
}

type groupMemberHandler struct {
	platform.BaseHandler[entities.GroupMember]
	service services.GroupMemberService
}
