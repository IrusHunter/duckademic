package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/services/auth/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// PermissionHandler represents a handler responsible for Permission-related HTTP operations.
type PermissionHandler interface {
	platform.BaseHandler[entities.Permission]
}

// NewPermissionHandler creates a new PermissionHandler instance.
//
// It requires a permission service.
func NewPermissionHandler(ps services.PermissionService) PermissionHandler {
	hc := platform.NewHandlerConfig("PermissionHandler", "permission")

	return &permissionHandler{
		BaseHandler: platform.NewBaseHandler(hc, ps),
		service:     ps,
	}
}

type permissionHandler struct {
	platform.BaseHandler[entities.Permission]
	service services.PermissionService
}
