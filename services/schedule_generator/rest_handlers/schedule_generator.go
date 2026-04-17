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
	SetLessonTypeAssignments(context.Context, http.ResponseWriter, *http.Request)
	SetStudentGroups(context.Context, http.ResponseWriter, *http.Request)
	SetStudyLoads(context.Context, http.ResponseWriter, *http.Request)
	SetClassrooms(context.Context, http.ResponseWriter, *http.Request)
	SubmitAndGoToTheNextStep(context.Context, http.ResponseWriter, *http.Request)
	SetDaysForLessonTypes(context.Context, http.ResponseWriter, *http.Request)
	GenerateBoneLessons(context.Context, http.ResponseWriter, *http.Request)
	AssignClassroomsToBoneLessons(context.Context, http.ResponseWriter, *http.Request)
	BuildScheduleSkeleton(context.Context, http.ResponseWriter, *http.Request)
	AddFloatingLessons(context.Context, http.ResponseWriter, *http.Request)
	AssignClassroomsToFloatingLessons(context.Context, http.ResponseWriter, *http.Request)
	GetStudyLoads(context.Context, http.ResponseWriter, *http.Request)
	GetLessons(context.Context, http.ResponseWriter, *http.Request)
	GetFault(context.Context, http.ResponseWriter, *http.Request)
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
	if h.generator == nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "CreateGenerator",
			fmt.Errorf("failed to set teachers: generator wasn't init"), logger.HandlerRequestFailed,
		))
		return
	}

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
	if h.generator == nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "CreateGenerator",
			fmt.Errorf("failed to set disciplines: generator wasn't init"), logger.HandlerRequestFailed,
		))
		return
	}

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
	if h.generator == nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "CreateGenerator",
			fmt.Errorf("failed to set lesson types: generator wasn't init"), logger.HandlerRequestFailed,
		))
		return
	}

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
func (h *scheduleGeneratorHandler) SetLessonTypeAssignments(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if h.generator == nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "CreateGenerator",
			fmt.Errorf("failed to set lesson type assignments: generator wasn't init"), logger.HandlerRequestFailed,
		))
		return
	}

	var assignments []entities.LessonTypeAssignment
	err := json.NewDecoder(r.Body).Decode(&assignments)
	defer r.Body.Close()
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"SetLessonTypeAssignments",
			fmt.Errorf("failed to extract lesson type assignments from body: %w", err),
			logger.HandlerRequestFailed,
		))
		return
	}

	if err := h.validationService.ValidateLessonTypeAssignments(assignments); err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"SetLessonTypeAssignments",
			fmt.Errorf("validation failed: %w", err),
			logger.HandlerRequestFailed,
		))
		return
	}

	if err := h.generator.SetLessonTypeAssignments(assignments); err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"SetLessonTypeAssignments",
			fmt.Errorf("failed to set lesson type assignments: %w", err),
			logger.HandlerRequestFailed,
		))
		return
	}

	jsonutil.ResponseWithJSON(w, 200, map[string]any{
		"message": fmt.Sprintf("%d lesson type assignments assigned", len(assignments)),
	})
}
func (h *scheduleGeneratorHandler) SetStudentGroups(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if h.generator == nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"SetStudentGroups",
			fmt.Errorf("failed to set group cohorts: generator wasn't init"),
			logger.HandlerRequestFailed,
		))
		return
	}

	req := struct {
		GroupCohorts           []entities.GroupCohort           `json:"group_cohorts"`
		GroupCohortAssignments []entities.GroupCohortAssignment `json:"group_cohort_assignments"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"SetStudentGroups",
			fmt.Errorf("failed to extract data from body: %w", err),
			logger.HandlerRequestFailed,
		))
		return
	}

	if err := h.validationService.ValidateGroupCohorts(req.GroupCohorts); err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"SetStudentGroups",
			fmt.Errorf("validation failed: %w", err),
			logger.HandlerRequestFailed,
		))
		return
	}

	if err := h.generator.SetStudentGroups(req.GroupCohorts, req.GroupCohortAssignments); err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"SetStudentGroups",
			fmt.Errorf("failed to set student groups: %w", err),
			logger.HandlerRequestFailed,
		))
		return
	}

	jsonutil.ResponseWithJSON(w, 200, map[string]any{
		"message": fmt.Sprintf("%d group cohorts assigned, %d assignments assigned",
			len(req.GroupCohorts), len(req.GroupCohortAssignments)),
	})
}
func (h *scheduleGeneratorHandler) SetStudyLoads(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if h.generator == nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SetStudyLoads",
			fmt.Errorf("failed to set teacher loads: generator wasn't init"), logger.HandlerRequestFailed,
		))
		return
	}

	var loads []entities.TeacherLoad
	err := json.NewDecoder(r.Body).Decode(&loads)
	defer r.Body.Close()
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SetStudyLoads",
			fmt.Errorf("failed to extract teacher loads from body: %w", err), logger.HandlerRequestFailed,
		))
		return
	}

	if err := h.generator.SetStudyLoads(loads); err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SetStudyLoads",
			fmt.Errorf("failed to set teacher loads: %w", err), logger.HandlerRequestFailed,
		))
		return
	}

	jsonutil.ResponseWithJSON(w, 200, map[string]any{
		"message": fmt.Sprintf("%d teacher loads assigned", len(loads)),
	})
}
func (h *scheduleGeneratorHandler) SetClassrooms(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if h.generator == nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"SetClassrooms",
			fmt.Errorf("failed to set classrooms: generator wasn't init"),
			logger.HandlerRequestFailed,
		))
		return
	}

	var classrooms []entities.Classroom
	err := json.NewDecoder(r.Body).Decode(&classrooms)
	defer r.Body.Close()
	if err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"SetClassrooms",
			fmt.Errorf("failed to extract classrooms from body: %w", err),
			logger.HandlerRequestFailed,
		))
		return
	}

	if err := h.validationService.ValidateClassrooms(classrooms); err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"SetClassrooms",
			fmt.Errorf("validation failed: %w", err),
			logger.HandlerRequestFailed,
		))
		return
	}

	if err := h.generator.SetClassrooms(classrooms); err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"SetClassrooms",
			fmt.Errorf("failed to set classrooms: %w", err),
			logger.HandlerRequestFailed,
		))
		return
	}

	jsonutil.ResponseWithJSON(
		w,
		200,
		map[string]any{"message": fmt.Sprintf("%d classrooms assigned", len(classrooms))},
	)
}
func (h *scheduleGeneratorHandler) SubmitAndGoToTheNextStep(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.execute(
		ctx,
		w,
		"submit_and_go_to_the_next_step",
		func() (any, error) {
			return h.generator.SubmitAndGoToTheNextStep()
		},
		"SetDaysForLessonTypes",
	)
}
func (h *scheduleGeneratorHandler) SetDaysForLessonTypes(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.execute(
		ctx,
		w,
		"set_days_for_lesson_types",
		func() (any, error) {
			return h.generator.SetDaysForLessonTypes()
		},
		"SetDaysForLessonTypes",
	)
}
func (h *scheduleGeneratorHandler) GenerateBoneLessons(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.execute(
		ctx,
		w,
		"generate_bone_lessons",
		func() (any, error) {
			return h.generator.GenerateBoneLessons()
		},
		"SetDaysForLessonTypes",
	)
}
func (h *scheduleGeneratorHandler) AssignClassroomsToBoneLessons(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.execute(
		ctx,
		w,
		"assign_classrooms_to_bone_lessons",
		func() (any, error) {
			return h.generator.AssignClassroomsToBoneLessons()
		},
		"AssignClassroomsToBoneLessons",
	)
}
func (h *scheduleGeneratorHandler) BuildScheduleSkeleton(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.execute(
		ctx,
		w,
		"build_lesson_skeleton",
		func() (any, error) {
			return h.generator.BuildScheduleSkeleton()
		},
		"BuildLessonSkeleton",
	)
}
func (h *scheduleGeneratorHandler) AddFloatingLessons(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.execute(
		ctx,
		w,
		"add_floating_lessons",
		func() (any, error) {
			return h.generator.AddFloatingLessons()
		},
		"AddFloatingLessons",
	)
}
func (h *scheduleGeneratorHandler) AssignClassroomsToFloatingLessons(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.execute(
		ctx,
		w,
		"assign_classrooms_to_floating_lessons",
		func() (any, error) {
			return h.generator.AssignClassroomsToFloatingLessons()
		},
		"AssignClassroomsToFloatingLessons",
	)
}
func (h *scheduleGeneratorHandler) GetStudyLoads(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.execute(
		ctx,
		w,
		"get_study_loads",
		func() (any, error) {
			return h.generator.ExtractStudyLoads()
		},
		"GetStudyLoads",
	)
}
func (h *scheduleGeneratorHandler) GetLessons(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.execute(
		ctx,
		w,
		"get_lessons",
		func() (any, error) {
			return h.generator.ExtractLessons()
		},
		"GetLessons",
	)
}
func (h *scheduleGeneratorHandler) GetFault(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.execute(
		ctx,
		w,
		"get_fault",
		func() (any, error) {
			return h.generator.GetFault(), nil
		},
		"GetFault",
	)
}

func (h *scheduleGeneratorHandler) execute(
	ctx context.Context,
	w http.ResponseWriter,
	action string,
	fn func() (any, error),
	name string,
) {
	if h.generator == nil {
		jsonutil.ResponseWithError(
			w,
			400,
			h.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				name,
				fmt.Errorf("generator wasn't initialized"),
				logger.HandlerRequestFailed,
			),
		)
		return
	}

	res, err := fn()
	if err != nil {
		jsonutil.ResponseWithError(
			w,
			400,
			h.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx),
				name,
				fmt.Errorf("generator action %s failed: %w", action, err),
				logger.HandlerRequestFailed,
			),
		)
		return
	}

	jsonutil.ResponseWithJSON(w, 200, res)
}
