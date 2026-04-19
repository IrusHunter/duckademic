package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/services/auth/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RolePermissionsHandler handles HTTP operations for role-permission assignments.
type RolePermissionsHandler interface {
	platform.BaseHandler[entities.RolePermissions]
}

// NewRolePermissionsHandler creates a new RolePermissionsHandler instance.
func NewRolePermissionsHandler(rps services.RolePermissionsService) RolePermissionsHandler {
	hc := platform.NewHandlerConfig("RolePermissionsHandler", "role_permission")

	return &rolePermissionsHandler{
		BaseHandler: platform.NewBaseHandler(hc, rps),
		service:     rps,
	}
}

type rolePermissionsHandler struct {
	platform.BaseHandler[entities.RolePermissions]
	service services.RolePermissionsService
}
