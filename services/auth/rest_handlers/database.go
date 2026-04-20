package resthandlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IrusHunter/duckademic/services/auth/services"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
)

type DatabaseHandler interface {
	Seed(context.Context, http.ResponseWriter, *http.Request)
	Clear(context.Context, http.ResponseWriter, *http.Request)
}

func NewDatabaseHandler(
	ps services.PermissionService,
	rs services.RoleService,
	rps services.RolePermissionsService,
	ss services.ServiceService,
	sps services.ServicePermissionsService,
	us services.UserService,
) DatabaseHandler {
	return &databaseHandler{
		permissionService:         ps,
		roleService:               rs,
		rolePermissionsService:    rps,
		servicePermissionsService: sps,
		serviceService:            ss,
		userService:               us,
	}
}

type databaseHandler struct {
	permissionService         services.PermissionService
	roleService               services.RoleService
	rolePermissionsService    services.RolePermissionsService
	serviceService            services.ServiceService
	servicePermissionsService services.ServicePermissionsService
	userService               services.UserService
}

func (h *databaseHandler) Seed(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.roleService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed roles: %w", err))
		return
	}
	if err := h.permissionService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed permissions: %w", err))
		return
	}
	if err := h.userService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed users: %w", err))
		return
	}
	if err := h.rolePermissionsService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed role permissions: %w", err))
		return
	}
	if err := h.serviceService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed services: %w", err))
		return
	}
	if err := h.servicePermissionsService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed service permissions: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
func (h *databaseHandler) Clear(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.servicePermissionsService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear service permissions: %w", err))
		return
	}
	if err := h.userService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear users: %w", err))
		return
	}
	if err := h.rolePermissionsService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear role permissions: %w", err))
		return
	}
	if err := h.roleService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear roles: %w", err))
		return
	}
	if err := h.serviceService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear services: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
