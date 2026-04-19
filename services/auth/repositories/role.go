package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type RoleRepository interface {
	platform.BaseRepository[entities.Role]
	FindByName(context.Context, string) *entities.Role
}

func NewRoleRepository(db *sqlx.DB) RoleRepository {
	config := platform.NewRepositoryConfig(
		"RoleRepository",
		entities.Role{}.TableName(),
		entities.Role{}.EntityName(),
		[]string{"id", "name"},
		[]string{"name"},
		[]string{"created_at", "updated_at"},
	)

	rr := &roleRepository{
		BaseRepository: platform.NewBaseRepository[entities.Role](config, db),
		db:             db,
	}
	rr.logger = rr.GetLogger()

	return rr
}

type roleRepository struct {
	platform.BaseRepository[entities.Role]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *roleRepository) FindByName(ctx context.Context, name string) *entities.Role {
	return r.FindFirstBy(ctx, "name", name)
}
