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
	SetTeachers(context.Context, http.ResponseWriter, *http.Request)
	SetDisciplines(context.Context, http.ResponseWriter, *http.Request)
	SetLessonTypes(context.Context, http.ResponseWriter, *http.Request)
}

func NewScheduleGeneratorHandler(
	gcs services.GeneratorConfigService,
	vs services.ValidationService,
) ScheduleGeneratorHandler {

	return &scheduleGeneratorHandler{
		generatorConfigService: gcs,
		validationService:      vs,
		logger:                 logger.NewLogger("ScheduleGeneratorHandler.txt", "ScheduleGeneratorHandler"),
	}
}

type scheduleGeneratorHandler struct {
	generatorConfigService services.GeneratorConfigService
	validationService      services.ValidationService
	generator              *core.ScheduleGenerator
	logger                 logger.Logger
}

func (h *scheduleGeneratorHandler) CreateGenerator(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if h.generator != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "CreateGenerator",
			fmt.Errorf("failed to create generator: generator already exists"), logger.HandlerRequestFailed,
		))
		return
	}

	var generatorConfig entities.ScheduleGeneratorConfig
	err := json.NewDecoder(r.Body).Decode(&generatorConfig)
	defer r.Body.Close()
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "CreateGenerator",
			fmt.Errorf("failed to extract generator config from body: %w", err), logger.HandlerRequestFailed,
		))
		return
	}

	if err := h.generatorConfigService.ValidateScheduleConfig(generatorConfig); err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "CreateGenerator",
			fmt.Errorf("validation failed: %w", err), logger.HandlerRequestFailed,
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

func (h *scheduleGeneratorHandler) SetTeachers(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var teachers []entities.Teacher
	err := json.NewDecoder(r.Body).Decode(&teachers)
	defer r.Body.Close()
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SetTeachers",
			fmt.Errorf("failed to extract teachers from body: %w", err), logger.HandlerRequestFailed,
		))
		return
	}

	if err := h.validationService.ValidateTeachers(teachers); err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SetTeachers",
			fmt.Errorf("validation failed: %w", err), logger.HandlerRequestFailed,
		))
		return
	}

	if err := h.generator.SetTeachers(teachers); err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SetTeachers",
			fmt.Errorf("failed to set teachers: %w", err), logger.HandlerRequestFailed,
		))
		return
	}

	jsonutil.ResponseWithJSON(w, 200, map[string]any{"message": fmt.Sprintf("%d teachers assigned", len(teachers))})
}
func (h *scheduleGeneratorHandler) SetDisciplines(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var disciplines []entities.Discipline
	err := json.NewDecoder(r.Body).Decode(&disciplines)
	defer r.Body.Close()
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SetDisciplines",
			fmt.Errorf("failed to extract disciplines from body: %w", err), logger.HandlerRequestFailed,
		))
		return
	}

	if err := h.validationService.ValidateDisciplines(disciplines); err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SetDisciplines",
			fmt.Errorf("validation failed: %w", err), logger.HandlerRequestFailed,
		))
		return
	}

	if err := h.generator.SetDisciplines(disciplines); err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SetDisciplines",
			fmt.Errorf("failed to set disciplines: %w", err), logger.HandlerRequestFailed,
		))
		return
	}

	jsonutil.ResponseWithJSON(w, 200, map[string]any{
		"message": fmt.Sprintf("%d disciplines assigned", len(disciplines)),
	})
}
func (h *scheduleGeneratorHandler) SetLessonTypes(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var requests []entities.LessonTypeRequest
	err := json.NewDecoder(r.Body).Decode(&requests)
	defer r.Body.Close()
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"SetLessonTypes",
			fmt.Errorf("failed to extract lesson types from body: %w", err),
			logger.HandlerRequestFailed,
		))
		return
	}

	if err := h.validationService.ValidateLessonTypeRequests(requests); err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"SetLessonTypes",
			fmt.Errorf("validation failed: %w", err),
			logger.HandlerRequestFailed,
		))
		return
	}

	lessonTypes := make([]entities.LessonType, 0, len(requests))
	for i, req := range requests {
		lt, err := req.ToLessonType()
		if err != nil {
			jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"SetLessonTypes",
				fmt.Errorf("conversion failed for item %d: %w", i, err),
				logger.HandlerRequestFailed,
			))
			return
		}

		lessonTypes = append(lessonTypes, lt)
	}

	if err := h.generator.SetLessonTypes(lessonTypes); err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"SetLessonTypes",
			fmt.Errorf("failed to set lesson types: %w", err),
			logger.HandlerRequestFailed,
		))
		return
	}

	jsonutil.ResponseWithJSON(w, 200, map[string]any{
		"message": fmt.Sprintf("%d lesson types assigned", len(lessonTypes)),
	})
}
