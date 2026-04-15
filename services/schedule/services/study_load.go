package services

import (
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/repositories"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type StudyLoadService interface {
	platform.BaseService[entities.StudyLoad]
}

func NewStudyLoadService(
	sr repositories.StudyLoadRepository,
) StudyLoadService {
	sc := platform.NewServiceConfig(
		"StudyLoadService",
		filepath.Join("data", "study_loads.json"),
		"study load",
	)

	res := &studyLoadService{
		repository: sr,
	}

	res.BaseService = platform.NewBaseService(
		sc,
		sr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.StudyLoad]{},
	)

	return res
}

type studyLoadService struct {
	platform.BaseService[entities.StudyLoad]
	repository repositories.StudyLoadRepository
}
