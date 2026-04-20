package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/services/auth/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type RolePermissionsService interface {
	platform.BaseService[entities.RolePermissions]
}

func NewRolePermissionsService(
	rpr repositories.RolePermissionsRepository,
	pr repositories.PermissionRepository,
	rr repositories.RoleRepository,
) RolePermissionsService {
	sc := platform.NewServiceConfig(
		"RolePermissionsService",
		filepath.Join("data", "role_permissions.json"),
		"role_permission",
	)

	s := &rolePermissionsService{
		repository:           rpr,
		permissionRepository: pr,
		roleRepository:       rr,
	}

	s.BaseService = platform.NewBaseService(sc, rpr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.RolePermissions]{
			platform.OnAddPrepare: s.onAddPrepare,
		},
	)

	s.logger = s.GetLogger()
	return s
}

type rolePermissionsService struct {
	platform.BaseService[entities.RolePermissions]
	repository           repositories.RolePermissionsRepository
	permissionRepository repositories.PermissionRepository
	roleRepository       repositories.RoleRepository
	logger               logger.Logger
}

func (s *rolePermissionsService) onAddPrepare(ctx context.Context, rp *entities.RolePermissions) error {
	rp.ID = uuid.New()
	return nil
}

// Seed loads role-permission assignments from JSON file and inserts them
func (s *rolePermissionsService) Seed(ctx context.Context) error {
	assignments := []struct {
		RoleName       string `json:"role_name"`
		PermissionName string `json:"permission_name"`
	}{}

	if err := jsonutil.ReadFileTo(filepath.Join("data", "role_permissions.json"), &assignments); err != nil {
		return s.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"Seed",
			fmt.Errorf("failed to load role permissions seed data: %w", err),
			logger.ServiceValidationFailed,
		)
	}

	var lastError error

	for _, item := range assignments {
		role := s.roleRepository.FindByName(ctx, item.RoleName)
		permission := s.permissionRepository.FindByName(ctx, item.PermissionName)

		if role == nil || permission == nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("role or permission not found: %s -> %s", item.RoleName, item.PermissionName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		rp := entities.RolePermissions{
			ID:           uuid.New(),
			RoleID:       role.ID,
			PermissionID: permission.ID,
		}

		_, err := s.Add(ctx, rp)
		if err != nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("failed to add role permission %s: %w", rp, err),
				logger.ServiceValidationFailed,
			)
			continue
		}
	}

	s.logger.Log(
		contextutil.GetTraceID(ctx),
		"Seed",
		"role permissions seed processed",
		logger.ServiceOperationSuccess,
	)

	return lastError
}
