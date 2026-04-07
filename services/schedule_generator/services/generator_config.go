package services

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/IrusHunter/duckademic/services/schedule_generator/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
)

type GeneratorConfigService interface {
	ValidateScheduleConfig(entities.ScheduleGeneratorConfig) error
	GetDefaultGeneratorConfig(context.Context) (entities.ScheduleGeneratorConfig, error)
}

func NewGeneratorConfigService() GeneratorConfigService {
	return &generatorConfigService{
		logger: logger.NewLogger("GeneratorConfigService.txt", "GeneratorConfigService"),
	}
}

type generatorConfigService struct {
	logger logger.Logger
}

func (s *generatorConfigService) ValidateScheduleConfig(cfg entities.ScheduleGeneratorConfig) error {
	if err := cfg.ValidateStartTime(); err != nil {
		return err
	}
	if err := cfg.ValidateEndTime(); err != nil {
		return err
	}
	if err := cfg.ValidateSlotPreference(); err != nil {
		return err
	}
	if err := cfg.ValidateMaxDailyStudentLoad(); err != nil {
		return err
	}
	if err := cfg.ValidateLessonFillRate(); err != nil {
		return err
	}
	if err := cfg.ValidateClassroomOccupancy(); err != nil {
		return err
	}

	return nil
}

func (s *generatorConfigService) GetDefaultGeneratorConfig(ctx context.Context) (entities.ScheduleGeneratorConfig, error) {
	cfg := entities.ScheduleGeneratorConfig{}

	if err := jsonutil.ReadFileTo(filepath.Join("data", "generator_config.json"), &cfg); err != nil {
		return entities.ScheduleGeneratorConfig{}, s.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to load default generator config: %w", err), logger.ServiceValidationFailed,
		)
	}

	return cfg, nil
}
