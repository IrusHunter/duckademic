package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RolePermissionsRepository interface {
	platform.BaseRepository[entities.RolePermissions]
	GetPermissionsByRoleID(context.Context, uuid.UUID) ([]string, error)
}

func NewRolePermissionsRepository(db *sqlx.DB) RolePermissionsRepository {
	config := platform.NewRepositoryConfig(
		"RolePermissionsRepository",
		entities.RolePermissions{}.TableName(),
		entities.RolePermissions{}.EntityName(),
		[]string{"id", "role_id", "permission_id"},
		[]string{},
		[]string{"created_at", "updated_at"},
	)

	r := &rolePermissionsRepository{
		BaseRepository: platform.NewBaseRepository[entities.RolePermissions](config, db),
		db:             db,
	}
	r.logger = r.GetLogger()

	return r
}

type rolePermissionsRepository struct {
	platform.BaseRepository[entities.RolePermissions]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *rolePermissionsRepository) GetPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]string, error) {
	query := `
		SELECT p.name
		FROM role_permissions rp
		JOIN permissions p ON rp.permission_id = p.id
		WHERE rp.role_id = $1;
	`

	var permissions []string

	if err := r.db.SelectContext(ctx, &permissions, query, roleID); err != nil {
		return nil, r.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetPermissionsByRoleID",
			err,
			logger.RepositoryScanFailed,
		)
	}

	return permissions, nil
}
