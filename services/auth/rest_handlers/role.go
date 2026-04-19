package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/services/auth/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RoleHandler represents a handler responsible for Role-related HTTP operations.
type RoleHandler interface {
	platform.BaseHandler[entities.Role]
}

// NewRoleHandler creates a new RoleHandler instance.
//
// It requires a role service.
func NewRoleHandler(rs services.RoleService) RoleHandler {
	hc := platform.NewHandlerConfig("RoleHandler", "role")

	return &roleHandler{
		BaseHandler: platform.NewBaseHandler(hc, rs),
		service:     rs,
	}
}

type roleHandler struct {
	platform.BaseHandler[entities.Role]
	service services.RoleService
}
