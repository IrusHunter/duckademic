package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type ServiceRepository interface {
	platform.BaseRepository[entities.Service]
	FindByName(ctx context.Context, name string) *entities.Service
}

func NewServiceRepository(db *sqlx.DB) ServiceRepository {
	config := platform.NewRepositoryConfig(
		"ServiceRepository",
		entities.Service{}.TableName(),
		entities.Service{}.EntityName(),
		[]string{"id", "name", "secrete"},
		[]string{"name", "secrete"},
		[]string{"created_at", "updated_at"},
	)

	r := &serviceRepository{
		BaseRepository: platform.NewBaseRepository[entities.Service](config, db),
		db:             db,
	}
	r.logger = r.GetLogger()

	return r
}

type serviceRepository struct {
	platform.BaseRepository[entities.Service]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *serviceRepository) FindByName(ctx context.Context, name string) *entities.Service {
	return r.FindFirstBy(ctx, "name", name)
}
