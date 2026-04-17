package services

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
)

type LessonOccurrenceService interface {
	platform.BaseService[entities.LessonOccurrence]
	AddFromExternal(ctx context.Context, el []entities.ExternalLesson) error
}

func NewLessonOccurrenceService(
	lr repositories.LessonOccurrenceRepository,
	lsr repositories.LessonSlotRepository,
) LessonOccurrenceService {
	sc := platform.NewServiceConfig(
		"LessonOccurrenceService",
		filepath.Join("data", "lesson_occurrences.json"),
		entities.LessonOccurrence{}.EntityName(),
	)

	res := &lessonOccurrenceService{
		repository:           lr,
		lessonSlotRepository: lsr,
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
	repository           repositories.LessonOccurrenceRepository
	lessonSlotRepository repositories.LessonSlotRepository
}

func (s *lessonOccurrenceService) AddFromExternal(ctx context.Context, el []entities.ExternalLesson) error {
	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var lastError error

	for i, externalL := range el {
		wg.Add(1)
		sem <- struct{}{}

		go func(i int, externalL entities.ExternalLesson) {
			defer wg.Done()
			defer func() { <-sem }()

			lesson := entities.LessonOccurrence{
				ID:             externalL.ID,
				StudyLoadID:    externalL.StudyLoadID,
				TeacherID:      externalL.TeacherID,
				StudentGroupID: externalL.StudentGroupID,
				ClassroomID:    externalL.ClassroomID,
				Status:         entities.LessonOccurrenceScheduled,
			}

			slot := s.lessonSlotRepository.FindBySlotAndWeekday(ctx, externalL.Slot, externalL.Day%7)
			if slot == nil {
				mu.Lock()
				lastError = s.GetLogger().LogAndReturnError(contextutil.GetTraceID(ctx), "AddMultiple",
					fmt.Errorf("failed to find lesson slot (%d/%d) [%d]", externalL.Day, externalL.Slot, i),
					logger.ServiceValidationFailed)
				mu.Unlock()
				return
			}

			lesson.LessonSlotID = slot.ID
			lesson.Date = time.Date(2026, time.January, 20, 0, 0, 0, 0, time.UTC).Add(slot.StartTime).
				Add(time.Hour * 24 * time.Duration(externalL.Day))

			_, err := s.Add(ctx, lesson)
			if err != nil {
				mu.Lock()
				lastError = s.GetLogger().LogAndReturnError(contextutil.GetTraceID(ctx), "AddMultiple",
					fmt.Errorf("failed to insert at index [%d]: %w", i, err), logger.ServiceValidationFailed)
				mu.Unlock()
			}
		}(i, externalL)

	}

	wg.Wait()
	return lastError
}
