package platform

import (
	"context"
	"fmt"

	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/google/uuid"
)

type ServiceConfig struct {
	ClassName  string
	JSONPath   string
	EntityName string
}

// NewServiceConfig creates a new ServiceConfig instance.
//
// It requires the class name (cn) path to the seed file (jsonP), and entity name (en).
func NewServiceConfig(cn, jsonP, en string) ServiceConfig {
	return ServiceConfig{
		ClassName:  cn,
		JSONPath:   jsonP,
		EntityName: en,
	}
}

type ServiceExternalFunc[T fmt.Stringer] func(context.Context, *T) error

type ServiceExternalFuncType int

const (
	OnAddPrepare ServiceExternalFuncType = iota
	ValidateEntity
	HardDeleteCheck
	OnUpdatePrepare
)

func getOrDefault[T fmt.Stringer](
	m map[ServiceExternalFuncType]ServiceExternalFunc[T],
	key ServiceExternalFuncType,
	def ServiceExternalFunc[T],
) ServiceExternalFunc[T] {
	if fn, ok := m[key]; ok {
		return fn
	}
	return def
}

// BaseService provides operations to initialize and manage entities.
type BaseService[T fmt.Stringer] interface {
	// Seed clears existing entities data and initializes it from a JSON file.
	Seed(context.Context) error
	Clear(context.Context) error
	// Add validates and inserts a new entity into the repository and returns it, or an error if it fails.
	Add(context.Context, T) (T, error)
	// GetAll returns a slice with all entities.
	GetAll(context.Context) []T
	// Delete removes the entity by ID, or marks it as deleted depending on service logic.
	Delete(context.Context, uuid.UUID) (T, error)
	// Update updates the entity with the specified ID and returns the updated one.
	Update(context.Context, uuid.UUID, T) (T, error)
	// FindByID returns a pointer to the entity from repository with the given id.
	FindByID(context.Context, uuid.UUID) *T
	GetLogger() logger.Logger
	SendChanges(context.Context, fmt.Stringer, events.EventType, events.RedisTopic)
}

// NewBaseService creates a new BaseService instance.
//
// It requires a repository (r), a config (sc), a validation (vf), and external functions (mf).
func NewBaseService[T fmt.Stringer](sc ServiceConfig, r BaseRepository[T],
	mf map[ServiceExternalFuncType]ServiceExternalFunc[T],
) BaseService[T] {
	onAddPrepare := getOrDefault(mf, OnAddPrepare, func(ctx context.Context, t *T) error { return nil })
	validateEntity := getOrDefault(mf, ValidateEntity, func(ctx context.Context, t *T) error { return nil })
	hardDeleteCheck := getOrDefault(mf, HardDeleteCheck, func(ctx context.Context, t *T) error { return nil })
	onUpdatePrepare := getOrDefault(mf, OnUpdatePrepare, func(ctx context.Context, t *T) error { return nil })

	return &baseService[T]{
		ServiceConfig:   sc,
		repository:      r,
		logger:          logger.NewLogger(sc.ClassName+".txt", sc.ClassName),
		validateEntity:  validateEntity,
		onAddPrepare:    onAddPrepare,
		hardDeleteCheck: hardDeleteCheck,
		onUpdatePrepare: onUpdatePrepare,
	}
}

func NewBaseServiceWithEventBus[T fmt.Stringer](sc ServiceConfig, r BaseRepository[T],
	mf map[ServiceExternalFuncType]ServiceExternalFunc[T], eb events.EventBus,
) BaseService[T] {
	onAddPrepare := getOrDefault(mf, OnAddPrepare, func(ctx context.Context, t *T) error { return nil })
	validateEntity := getOrDefault(mf, ValidateEntity, func(ctx context.Context, t *T) error { return nil })
	hardDeleteCheck := getOrDefault(mf, HardDeleteCheck, func(ctx context.Context, t *T) error { return nil })
	onUpdatePrepare := getOrDefault(mf, OnUpdatePrepare, func(ctx context.Context, t *T) error { return nil })

	return &baseService[T]{
		ServiceConfig:   sc,
		repository:      r,
		logger:          logger.NewLogger(sc.ClassName+".txt", sc.ClassName),
		validateEntity:  validateEntity,
		onAddPrepare:    onAddPrepare,
		hardDeleteCheck: hardDeleteCheck,
		onUpdatePrepare: onUpdatePrepare,
		eventBus:        eb,
	}
}

type baseService[T fmt.Stringer] struct {
	ServiceConfig
	repository      BaseRepository[T]
	logger          logger.Logger
	validateEntity  func(context.Context, *T) error
	onAddPrepare    func(context.Context, *T) error
	onUpdatePrepare func(context.Context, *T) error
	hardDeleteCheck func(context.Context, *T) error
	nilEntity       T
	eventBus        events.EventBus
}

func (s *baseService[T]) Add(ctx context.Context, entity T) (T, error) {
	if err := s.validateEntity(ctx, &entity); err != nil {
		return s.nilEntity, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Add",
			fmt.Errorf("%s failed validation: %w", entity.String(), err), logger.ServiceValidationFailed,
		)
	}

	if err := s.onAddPrepare(ctx, &entity); err != nil {
		return s.nilEntity, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Add",
			fmt.Errorf("failed to prepare entity %s: %w", entity.String(), err), logger.ServiceValidationFailed,
		)
	}

	ar, err := s.repository.Add(ctx, entity)
	if err != nil {
		return s.nilEntity, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Add",
			err, logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Add",
		fmt.Sprintf("%s added successfully", ar), logger.ServiceOperationSuccess,
	)
	return ar, nil
}
func (s *baseService[T]) Seed(ctx context.Context) error {
	entities := []T{}
	if err := jsonutil.ReadFileTo(s.JSONPath, &entities); err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load %ss seed data: %w", s.EntityName, err), logger.ServiceValidationFailed,
		)
	}

	var lastError error
	for _, entity := range entities {
		_, err := s.Add(ctx, entity)
		if err != nil {
			lastError = s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
				fmt.Errorf("failed to add %s: %w", entity.String(), err), logger.ServiceValidationFailed,
			)
		}
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Seed",
		fmt.Sprintf("%d %ss added successfully", len(entities), s.EntityName), logger.ServiceOperationSuccess,
	)
	return lastError
}
func (s *baseService[T]) Clear(ctx context.Context) error {
	err := s.repository.Clear(ctx)
	if err != nil {
		return s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Clear",
			err, logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Clear",
		fmt.Sprintf("all %ss cleared", s.EntityName), logger.ServiceOperationSuccess)
	return nil
}
func (s *baseService[T]) GetAll(ctx context.Context) []T {
	res := s.repository.GetAll(ctx)

	s.logger.Log(contextutil.GetTraceID(ctx), "GetAll",
		fmt.Sprintf("%d entities found", len(res)), logger.ServiceOperationSuccess)
	return res
}
func (s *baseService[T]) Delete(ctx context.Context, id uuid.UUID) (T, error) {
	entity := s.FindByID(ctx, id)
	if entity == nil {
		return s.nilEntity, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Delete",
			fmt.Errorf("%s with id %q not found", s.EntityName, id), logger.ServiceValidationFailed,
		)
	}

	if err := s.hardDeleteCheck(ctx, entity); err != nil {
		deletedE, err := s.repository.SoftDelete(ctx, id)
		if err != nil {
			return s.nilEntity, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SoftDelete",
				fmt.Errorf("failed to soft delete %s with id %q in repository: %w", s.EntityName, id, err),
				logger.ServiceRepositoryFailed,
			)
		}

		s.logger.Log(contextutil.GetTraceID(ctx), "SoftDelete",
			fmt.Sprintf("%s successfully soft deleted", deletedE),
			logger.ServiceOperationSuccess)
		return deletedE, nil
	} else {
		err := s.repository.Delete(ctx, id)
		if err != nil {
			return s.nilEntity, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Delete",
				fmt.Errorf("failed to hard delete %s with id %q from repository: %w", s.EntityName, id, err),
				logger.ServiceRepositoryFailed,
			)
		}

		s.logger.Log(contextutil.GetTraceID(ctx), "Delete",
			fmt.Sprintf("%s successfully hard deleted", *entity), logger.ServiceOperationSuccess)
		return *entity, nil
	}
}
func (s *baseService[T]) Update(ctx context.Context, id uuid.UUID, entity T) (T, error) {
	if err := s.validateEntity(ctx, &entity); err != nil {
		return s.nilEntity, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Update",
			fmt.Errorf("%s failed validation: %w", entity, err), logger.ServiceValidationFailed,
		)
	}

	updatedE, err := s.repository.Update(ctx, id, entity)
	if err != nil {
		return s.nilEntity, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Update",
			err, logger.ServiceRepositoryFailed,
		)
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "Update",
		fmt.Sprintf("%s successfully updated", updatedE), logger.ServiceOperationSuccess)
	return updatedE, nil
}
func (s *baseService[T]) FindByID(ctx context.Context, id uuid.UUID) *T {
	entity := s.repository.FindByID(ctx, id)
	if entity == nil {
		s.logger.Log(contextutil.GetTraceID(ctx), "FindByID",
			fmt.Sprintf("%s with id %q not found", s.EntityName, id), logger.ServiceOperationSuccess,
		)
	} else {
		s.logger.Log(contextutil.GetTraceID(ctx), "FindByID",
			fmt.Sprintf("%s found", entity), logger.ServiceOperationSuccess,
		)
	}
	return entity
}
func (s *baseService[T]) GetLogger() logger.Logger {
	return s.logger
}
func (s *baseService[T]) SendChanges(
	ctx context.Context, entity fmt.Stringer, event events.EventType, topic events.RedisTopic,
) {
	bytesEAR, err := events.ToByteConvertor(entity)
	if err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SendChanges",
			fmt.Errorf("failed to convert %s to bytes: %w", entity, err), logger.EventDataReadFailed,
		)
		return
	}

	if err := s.eventBus.Publish(ctx, string(topic), bytesEAR); err != nil {
		s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SendChanges",
			fmt.Errorf("failed to publish %s: %w", entity, err), logger.EventDataSentFailed,
		)
		return
	}

	s.logger.Log(contextutil.GetTraceID(ctx), "SendChanges",
		fmt.Sprintf("%s published successfully", entity), logger.EventDataSentSuccess)
}
