package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/services/auth/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// ServicePermissionsHandler handles HTTP operations for service-permission assignments.
type ServicePermissionsHandler interface {
	platform.BaseHandler[entities.ServicePermissions]
}

// NewServicePermissionsHandler creates a new ServicePermissionsHandler instance.
func NewServicePermissionsHandler(sps services.ServicePermissionsService) ServicePermissionsHandler {
	hc := platform.NewHandlerConfig("ServicePermissionsHandler", "service_permission")

	return &servicePermissionsHandler{
		BaseHandler: platform.NewBaseHandler(hc, sps),
		service:     sps,
	}
}

type servicePermissionsHandler struct {
	platform.BaseHandler[entities.ServicePermissions]
	service services.ServicePermissionsService
}
