package services

import (
	"context"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/services/auth/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type RoleService interface {
	platform.BaseService[entities.Role]
}

func NewRoleService(rr repositories.RoleRepository, ar string) (RoleService, uuid.UUID) {
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

	s.Add(contextutil.SetTraceID(context.Background()), entities.Role{Name: ar})
	s.Add(contextutil.SetTraceID(context.Background()), entities.Role{Name: "teacher"})
	s.Add(contextutil.SetTraceID(context.Background()), entities.Role{Name: "student"})

	adminRole := s.repository.FindByName(contextutil.SetTraceID(context.Background()), ar)

	if adminRole == nil {
		adminRole = &entities.Role{}
	}

	return s, adminRole.ID
}

type roleService struct {
	platform.BaseService[entities.Role]
	repository repositories.RoleRepository
}

func (s *roleService) onAddPrepare(ctx context.Context, r *entities.Role) error {
	r.ID = uuid.New()
	return nil
}
