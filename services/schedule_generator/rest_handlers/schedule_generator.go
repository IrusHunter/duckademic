package resthandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core"
	"github.com/IrusHunter/duckademic/services/schedule_generator/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/services"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
)

type ScheduleGeneratorHandler interface {
	CreateGenerator(context.Context, http.ResponseWriter, *http.Request)
	GetDefaultConfig(context.Context, http.ResponseWriter, *http.Request)
}

func NewScheduleGeneratorHandler(gcs services.GeneratorConfigService) ScheduleGeneratorHandler {

	return &scheduleGeneratorHandler{
		generatorConfigService: gcs,
		logger:                 logger.NewLogger("ScheduleGeneratorHandler.txt", "ScheduleGeneratorHandler"),
	}
}

type scheduleGeneratorHandler struct {
	generatorConfigService services.GeneratorConfigService
	generator              *core.ScheduleGenerator
	logger                 logger.Logger
}

func (h *scheduleGeneratorHandler) CreateGenerator(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if h.generator != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "CreateGenerator",
			fmt.Errorf("failed to create generator: generator already exists"), logger.HandlerBadRequest,
		))
		return
	}

	var generatorConfig entities.ScheduleGeneratorConfig
	err := json.NewDecoder(r.Body).Decode(&generatorConfig)
	defer r.Body.Close()
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "CreateGenerator",
			fmt.Errorf("failed to extract generator config from body: %w", err), logger.HandlerBadRequest,
		))
		return
	}

	if err := h.generatorConfigService.ValidateScheduleConfig(generatorConfig); err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "CreateGenerator",
			fmt.Errorf("validation failed: %w", err), logger.HandlerBadRequest,
		))
		return
	}

	h.generator, err = core.NewScheduleGenerator(generatorConfig)
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "CreateGenerator",
			fmt.Errorf("failed to create generator: %w", err), logger.HandlerBadRequest,
		))
		return
	}

	jsonutil.ResponseWithJSON(w, 201, nil)
}

func (h *scheduleGeneratorHandler) GetDefaultConfig(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	cfg, err := h.generatorConfigService.GetDefaultGeneratorConfig(ctx)
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "GetDefaultConfig",
			err, logger.HandlerInternalError))
		return
	}

	jsonutil.ResponseWithJSON(w, 200, cfg)
}
