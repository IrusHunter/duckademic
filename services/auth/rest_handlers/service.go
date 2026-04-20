package resthandlers

import (
	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/services/auth/services"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// ServiceHandler handles HTTP operations for Service entities.
type ServiceHandler interface {
	platform.BaseHandler[entities.Service]
}

// NewServiceHandler creates a new ServiceHandler instance.
func NewServiceHandler(ss services.ServiceService) ServiceHandler {
	hc := platform.NewHandlerConfig("ServiceHandler", "service")

	return &serviceHandler{
		BaseHandler: platform.NewBaseHandler(hc, ss),
		service:     ss,
	}
}

type serviceHandler struct {
	platform.BaseHandler[entities.Service]
	service services.ServiceService
}
