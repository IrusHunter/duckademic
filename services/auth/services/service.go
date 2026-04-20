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

type ServiceService interface {
	platform.BaseService[entities.Service]
}

func NewServiceService(sr repositories.ServiceRepository) ServiceService {
	sc := platform.NewServiceConfig(
		"ServiceService",
		filepath.Join("data", "services.json"),
		"service",
	)

	s := &serviceService{
		repository: sr,
	}

	s.BaseService = platform.NewBaseService(
		sc,
		sr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.Service]{
			platform.OnAddPrepare: s.onAddPrepare,
		},
	)

	s.logger = s.GetLogger()

	return s
}

type serviceService struct {
	platform.BaseService[entities.Service]
	repository repositories.ServiceRepository
	logger     logger.Logger
}

func (s *serviceService) onAddPrepare(ctx context.Context, srv *entities.Service) error {
	srv.ID = uuid.New()
	return nil
}
