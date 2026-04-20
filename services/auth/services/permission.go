package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/services/auth/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/envutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type PermissionService interface {
	platform.BaseService[entities.Permission]
}

func NewPermissionService(
	pr repositories.PermissionRepository,
	rr repositories.RoleRepository,
	rpr repositories.RolePermissionsRepository,
	eb events.EventBus,
	arID uuid.UUID,
) PermissionService {
	sc := platform.NewServiceConfig(
		"PermissionService",
		filepath.Join("data", "permissions.json"),
		"permission",
	)

	s := &permissionService{
		repository:                pr,
		roleRepository:            rr,
		rolePermissionsRepository: rpr,
		adminRoleID:               arID,
	}

	s.BaseService = platform.NewBaseService(
		sc,
		pr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Permission]{
			platform.OnAddPrepare: s.onAddPrepare,
		},
	)

	s.logger = s.GetLogger()

	eb.Subscribe(
		contextutil.SetTraceID(context.Background()),
		string(events.AccessPermissionRT),
		s.eventHandler,
	)

	return s
}

type permissionService struct {
	platform.BaseService[entities.Permission]
	repository                repositories.PermissionRepository
	roleRepository            repositories.RoleRepository
	rolePermissionsRepository repositories.RolePermissionsRepository
	logger                    logger.Logger
	adminRoleID               uuid.UUID
}

func (s *permissionService) eventHandler(ctx context.Context, b []byte) {
	pr, err := events.FromByteConvertor[events.AccessPermissionRE](b)
	if err != nil {
		s.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"PermissionRTHandler",
			err,
			logger.EventDataReadFailed,
		)
		return
	}

	s.logger.Log(
		contextutil.GetTraceID(ctx),
		"PermissionRTHandler",
		fmt.Sprintf("received %s", pr),
		logger.EventDataReceived,
	)

	truePR := entities.Permission{
		Name: pr.Name,
	}

	s.Add(ctx, truePR)
}
func (s *permissionService) onAddPrepare(ctx context.Context, p *entities.Permission) error {
	p.ID = uuid.New()

	return nil
}

func (s *permissionService) Add(ctx context.Context, p entities.Permission) (entities.Permission, error) {
	res, err := s.BaseService.Add(ctx, p)
	if err == nil {
		s.rolePermissionsRepository.Add(ctx, entities.RolePermissions{
			RoleID:       s.adminRoleID,
			PermissionID: res.ID,
		})
	}
	return res, err
}

func (s *permissionService) Seed(ctx context.Context) error {
	adminRoleStr := envutil.GetStringFromENV("SUPER_ADMIN_ROLE")
	if adminRoleStr == "" {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("SUPER_ADMIN_ROLE not specified in the .env file"), logger.ServiceDataFetchFailed)
	}

	adminRole := s.roleRepository.FindByName(ctx, adminRoleStr)

	if adminRole == nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to sind admin role: %s", adminRole), logger.ServiceDataFetchFailed)
	}

	s.adminRoleID = adminRole.ID
	return nil
}

// func (s *permissionService) ExternalUpdate(
// 	ctx context.Context,
// 	id uuid.UUID,
// 	permission entities.Permission,
// ) (entities.Permission, error) {
// 	updated, err := s.repository.ExternalUpdate(ctx, id, permission)
// 	if err != nil {
// 		return entities.Permission{}, s.logger.LogAndReturnError(
// 			contextutil.GetTraceID(ctx),
// 			"ExternalUpdate",
// 			err,
// 			logger.ServiceRepositoryFailed,
// 		)
// 	}

// 	s.logger.Log(
// 		contextutil.GetTraceID(ctx),
// 		"ExternalUpdate",
// 		fmt.Sprintf("%s successfully updated", updated),
// 		logger.ServiceOperationSuccess,
// 	)

// 	return updated, nil
// }
