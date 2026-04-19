package repositories

import (
	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type RolePermissionsRepository interface {
	platform.BaseRepository[entities.RolePermissions]
	// FindByRoleID(ctx context.Context, roleID uuid.UUID) []entities.RolePermissions
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

// func (r *rolePermissionsRepository) FindByPermissionID(ctx context.Context, permissionID uuid.UUID) []entities.RolePermissions {
// 	return r.FindAllBy(ctx, "permission_id", permissionID)
//}
