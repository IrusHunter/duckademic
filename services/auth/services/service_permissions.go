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

type ServicePermissionsService interface {
	platform.BaseService[entities.ServicePermissions]
}

func NewServicePermissionsService(
	spr repositories.ServicePermissionsRepository,
	sr repositories.ServiceRepository,
	pr repositories.PermissionRepository,
) ServicePermissionsService {
	sc := platform.NewServiceConfig(
		"ServicePermissionsService",
		filepath.Join("data", "service_permissions.json"),
		"service_permission",
	)

	s := &servicePermissionsService{
		repository:           spr,
		serviceRepository:    sr,
		permissionRepository: pr,
	}

	s.BaseService = platform.NewBaseService(sc, spr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.ServicePermissions]{
			platform.OnAddPrepare: s.onAddPrepare,
		},
	)

	s.logger = s.GetLogger()
	return s
}

type servicePermissionsService struct {
	platform.BaseService[entities.ServicePermissions]
	repository           repositories.ServicePermissionsRepository
	serviceRepository    repositories.ServiceRepository
	permissionRepository repositories.PermissionRepository
	logger               logger.Logger
}

func (s *servicePermissionsService) onAddPrepare(ctx context.Context, sp *entities.ServicePermissions) error {
	sp.ID = uuid.New()
	return nil
}

// Seed loads service-permission assignments from JSON file and inserts them
func (s *servicePermissionsService) Seed(ctx context.Context) error {
	assignments := []struct {
		ServiceName    string `json:"service_name"`
		PermissionName string `json:"permission_name"`
	}{}

	if err := jsonutil.ReadFileTo(filepath.Join("data", "service_permissions.json"), &assignments); err != nil {
		return s.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"Seed",
			fmt.Errorf("failed to load service permissions seed data: %w", err),
			logger.ServiceValidationFailed,
		)
	}

	var lastError error

	for _, item := range assignments {
		service := s.serviceRepository.FindByName(ctx, item.ServiceName)
		permission := s.permissionRepository.FindByName(ctx, item.PermissionName)

		if service == nil || permission == nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("service or permission not found: %s -> %s", item.ServiceName, item.PermissionName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		sp := entities.ServicePermissions{
			ID:           uuid.New(),
			ServiceID:    service.ID,
			PermissionID: permission.ID,
		}

		_, err := s.Add(ctx, sp)
		if err != nil {
			lastError = s.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("failed to add service permission %s: %w", sp, err),
				logger.ServiceValidationFailed,
			)
			continue
		}
	}

	s.logger.Log(
		contextutil.GetTraceID(ctx),
		"Seed",
		"service permissions seed processed",
		logger.ServiceOperationSuccess,
	)

	return lastError
}
