package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PermissionRepository interface {
	platform.BaseRepository[entities.Permission]
	FindByName(context.Context, string) *entities.Permission
	ExternalUpdate(context.Context, uuid.UUID, entities.Permission) (entities.Permission, error)
}

func NewPermissionRepository(db *sqlx.DB) PermissionRepository {
	config := platform.NewRepositoryConfig(
		"PermissionRepository",
		entities.Permission{}.TableName(),
		entities.Permission{}.EntityName(),
		[]string{"id", "name"},
		[]string{"name"},
		[]string{"created_at", "updated_at"},
	)

	pr := &permissionRepository{
		BaseRepository: platform.NewBaseRepository[entities.Permission](config, db),
		db:             db,
	}
	pr.logger = pr.GetLogger()

	return pr
}

type permissionRepository struct {
	platform.BaseRepository[entities.Permission]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *permissionRepository) FindByName(ctx context.Context, name string) *entities.Permission {
	return r.FindFirstBy(ctx, "name", name)
}

func (r *permissionRepository) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	permission entities.Permission,
) (entities.Permission, error) {
	return r.UpdateFields(ctx, id, []string{"name"}, permission)
}
