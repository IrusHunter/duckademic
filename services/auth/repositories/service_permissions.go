package repositories

import (
	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type ServicePermissionsRepository interface {
	platform.BaseRepository[entities.ServicePermissions]
}

func NewServicePermissionsRepository(db *sqlx.DB) ServicePermissionsRepository {
	config := platform.NewRepositoryConfig(
		"ServicePermissionsRepository",
		entities.ServicePermissions{}.TableName(),
		entities.ServicePermissions{}.EntityName(),
		[]string{"id", "service_id", "permission_id"},
		[]string{},
		[]string{"created_at", "updated_at"},
	)

	r := &servicePermissionsRepository{
		BaseRepository: platform.NewBaseRepository[entities.ServicePermissions](config, db),
		db:             db,
	}
	r.logger = r.GetLogger()

	return r
}

type servicePermissionsRepository struct {
	platform.BaseRepository[entities.ServicePermissions]
	db     *sqlx.DB
	logger logger.Logger
}
