package services

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

type TeacherLoadService interface {
	platform.BaseService[entities.TeacherLoad]
	GetByLessonTypeAssignments(context.Context, []entities.LessonTypeAssignment) ([]entities.TeacherLoad, error)
	GetUniqueTeacherIDs([]entities.TeacherLoad) []uuid.UUID
	ToGeneratorTeacherLoads(context.Context, []entities.TeacherLoad) []GeneratorTeacherLoad
}

func NewTeacherLoadService(tr repositories.TeacherLoadRepository, eb events.EventBus) TeacherLoadService {
	sc := platform.NewServiceConfig("TeacherLoadService", filepath.Join("data", "teacher_loads.json"), "teacher_load")

	res := &teacherLoadService{
		repository: tr,
	}
	res.BaseService = platform.NewBaseService(sc, tr,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.TeacherLoad]{},
	)
	res.logger = res.GetLogger()

	eb.Subscribe(contextutil.SetTraceID(context.Background()), string(events.TeacherLoadRT), res.eventHandler)

	return res
}

type teacherLoadService struct {
	platform.BaseService[entities.TeacherLoad]
	repository repositories.TeacherLoadRepository
	logger     logger.Logger
}

func (s *teacherLoadService) eventHandler(ctx context.Context, b []byte) {
	loadEvent, err := events.FromByteConvertor[events.TeacherLoadRE](b)
	if err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "TeacherLoadRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "TeacherLoadRTHandler",
		fmt.Sprintf("received %s", loadEvent), logger.EventDataReceived,
	)

	trueLoad := entities.TeacherLoad{
		ID:           loadEvent.ID,
		TeacherID:    loadEvent.TeacherID,
		DisciplineID: loadEvent.DisciplineID,
		LessonTypeID: loadEvent.LessonTypeID,
		GroupCount:   loadEvent.GroupCount,
	}

	switch loadEvent.Event {
	case events.EntityCreated:
		s.Add(ctx, trueLoad)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, loadEvent.ID, trueLoad)
	case events.EntityDeleted:
		s.Delete(ctx, loadEvent.ID)
	}
}

func (s *teacherLoadService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	load entities.TeacherLoad,
) (entities.TeacherLoad, error) {
	updatedLoad, err := s.repository.ExternalUpdate(ctx, id, load)
	if err != nil {
		return entities.TeacherLoad{}, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "ExternalUpdate",
			err, logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedLoad), logger.ServiceOperationSuccess)
	return updatedLoad, nil
}

func (s *teacherLoadService) GetByLessonTypeAssignments(
	ctx context.Context,
	lta []entities.LessonTypeAssignment,
) ([]entities.TeacherLoad, error) {
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup
	var mu sync.Mutex

	var result []entities.TeacherLoad
	var lastError error

	for _, lessonTypeAssignment := range lta {
		wg.Add(1)
		sem <- struct{}{}

		go func(lessonTypeAssignment entities.LessonTypeAssignment) {
			defer wg.Done()
			defer func() { <-sem }()

			loads, err := s.repository.GetByDisciplineAndLessonTypeID(
				ctx, lessonTypeAssignment.DisciplineID, lessonTypeAssignment.LessonTypeID)
			if err != nil {
				mu.Lock()
				lastError = s.GetLogger().LogAndReturnError(
					contextutil.GetTraceID(ctx),
					"GetByDisciplineIDs",
					err,
					logger.ServiceRepositoryFailed,
				)
				mu.Unlock()
				return
			}

			mu.Lock()
			result = append(result, loads...)
			mu.Unlock()

		}(lessonTypeAssignment)
	}

	wg.Wait()

	return result, lastError
}

func (s *teacherLoadService) GetUniqueTeacherIDs(loads []entities.TeacherLoad) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{})
	result := make([]uuid.UUID, 0)

	for _, item := range loads {
		if _, ok := seen[item.TeacherID]; !ok {
			seen[item.TeacherID] = struct{}{}
			result = append(result, item.TeacherID)
		}
	}

	return result
}

type GeneratorTeacherLoad struct {
	ID           uuid.UUID `db:"id" json:"id"`
	TeacherID    uuid.UUID `db:"teacher_id" json:"teacher_id"`
	DisciplineID uuid.UUID `db:"discipline_id" json:"discipline_id"`
	LessonTypeID uuid.UUID `db:"lesson_type_id" json:"lesson_type_id"`
	GroupCount   int       `db:"group_count" json:"group_count"`
}

func (s *teacherLoadService) ToGeneratorTeacherLoads(ctx context.Context, tl []entities.TeacherLoad) []GeneratorTeacherLoad {
	res := make([]GeneratorTeacherLoad, 0, len(tl))

	for _, teacherLoad := range tl {
		res = append(res, GeneratorTeacherLoad{
			ID:           teacherLoad.ID,
			TeacherID:    teacherLoad.TeacherID,
			DisciplineID: teacherLoad.DisciplineID,
			LessonTypeID: teacherLoad.LessonTypeID,
			GroupCount:   teacherLoad.GroupCount,
		})
	}

	return res
}
