package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/services/auth/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type PermissionService interface {
	platform.BaseService[entities.Permission]
}

func NewPermissionService(pr repositories.PermissionRepository, eb events.EventBus) PermissionService {
	sc := platform.NewServiceConfig(
		"PermissionService",
		filepath.Join("data", "permissions.json"),
		"permission",
	)

	s := &permissionService{
		repository: pr,
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
	repository repositories.PermissionRepository
	logger     logger.Logger
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
