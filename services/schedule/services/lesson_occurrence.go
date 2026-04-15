package services

import (
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/repositories"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type LessonOccurrenceService interface {
	platform.BaseService[entities.LessonOccurrence]
}

func NewLessonOccurrenceService(
	lr repositories.LessonOccurrenceRepository,
) LessonOccurrenceService {
	sc := platform.NewServiceConfig(
		"LessonOccurrenceService",
		filepath.Join("data", "lesson_occurrences.json"),
		entities.LessonOccurrence{}.EntityName(),
	)

	res := &lessonOccurrenceService{
		repository: lr,
	}

	res.BaseService = platform.NewBaseService(
		sc,
		lr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.LessonOccurrence]{},
	)

	return res
}

type lessonOccurrenceService struct {
	platform.BaseService[entities.LessonOccurrence]
	repository repositories.LessonOccurrenceRepository
}
