package services

import (
	"context"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/services/auth/repositories"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type RoleService interface {
	platform.BaseService[entities.Role]
}

func NewRoleService(rr repositories.RoleRepository) RoleService {
	sc := platform.NewServiceConfig(
		"RoleService",
		filepath.Join("data", "roles.json"),
		"role",
	)

	s := &roleService{
		repository: rr,
	}

	s.BaseService = platform.NewBaseService(
		sc,
		rr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Role]{
			platform.OnAddPrepare: s.onAddPrepare,
		},
	)

	s.logger = s.GetLogger()

	return s
}

type roleService struct {
	platform.BaseService[entities.Role]
	repository repositories.RoleRepository
	logger     logger.Logger
}

func (s *roleService) onAddPrepare(ctx context.Context, r *entities.Role) error {
	r.ID = uuid.New()
	return nil
}
