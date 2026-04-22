package resthandlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/google/uuid"
)

// DatabaseHandler represents a handler responsible for database-related HTTP operations.
type DatabaseHandler interface {
	// Performs database seeding operations, initializing required data.
	Seed(context.Context, http.ResponseWriter, *http.Request)
	Clear(context.Context, http.ResponseWriter, *http.Request)
	ExtractDataFromGenerator(context.Context, http.ResponseWriter, *http.Request)
	LoadDataIntoGenerator(context.Context, http.ResponseWriter, *http.Request)
}

func NewDatabaseHandler(
	httpC *http.Client,
	sgd string,
	ars services.AcademicRankService,
	ts services.TeacherService,
	ds services.DisciplineService,
	lts services.LessonTypeService,
	ltas services.LessonTypeAssignmentService,
	ss services.StudentService,
	sgs services.StudentGroupService,
	gmg services.GroupMemberService,
	tls services.TeacherLoadService,
	gcs services.GroupCohortService,
	gcas services.GroupCohortAssignmentService,
	cs services.ClassroomService,
	sls services.StudyLoadService,
	lsl services.LessonSlotService,
	los services.LessonOccurrenceService,
	semS services.SemesterService,
	sds services.SemesterDisciplineService,
) DatabaseHandler {
	return &databaseHandler{
		httpClient:                   httpC,
		scheduleGeneratorDomain:      sgd,
		logger:                       logger.NewLogger("DatabaseHandler.txt", "DatabaseHandler"),
		academicRankService:          ars,
		teacherService:               ts,
		disciplineService:            ds,
		lessonTypeService:            lts,
		lessonTypeAssignmentService:  ltas,
		studentService:               ss,
		studentGroupService:          sgs,
		groupMemberService:           gmg,
		teacherLoadService:           tls,
		groupCohortService:           gcs,
		groupCohortAssignmentService: gcas,
		classroomService:             cs,
		studyLoadService:             sls,
		lessonSlotService:            lsl,
		lessonOccurrenceService:      los,
		semesterService:              semS,
		semesterDisciplineService:    sds,
	}
}

type databaseHandler struct {
	httpClient                   *http.Client
	scheduleGeneratorDomain      string
	logger                       logger.Logger
	academicRankService          services.AcademicRankService
	teacherService               services.TeacherService
	disciplineService            services.DisciplineService
	lessonTypeService            services.LessonTypeService
	lessonTypeAssignmentService  services.LessonTypeAssignmentService
	studentService               services.StudentService
	studentGroupService          services.StudentGroupService
	groupMemberService           services.GroupMemberService
	teacherLoadService           services.TeacherLoadService
	groupCohortService           services.GroupCohortService
	groupCohortAssignmentService services.GroupCohortAssignmentService
	classroomService             services.ClassroomService
	studyLoadService             services.StudyLoadService
	lessonSlotService            services.LessonSlotService
	lessonOccurrenceService      services.LessonOccurrenceService
	semesterService              services.SemesterService
	semesterDisciplineService    services.SemesterDisciplineService
}

func (h *databaseHandler) Seed(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.lessonSlotService.Seed(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed lesson slots: %w", err))
		return
	}

	go func() {
		time.Sleep(events.ExternalSeedCooldown)
		ctx := contextutil.SetTraceID(context.Background())
		h.academicRankService.Seed(ctx)
		ctx = contextutil.SetTraceID(context.Background())
		h.lessonTypeService.Seed(ctx)
	}()

	jsonutil.ResponseWithJSON(w, 204, nil)
}
func (h *databaseHandler) Clear(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.semesterDisciplineService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear semester disciplines: %w", err))
		return
	}
	if err := h.classroomService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear classrooms: %w", err))
		return
	}
	if err := h.teacherLoadService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear teacher loads: %w", err))
		return
	}
	if err := h.groupMemberService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear group members: %w", err))
		return
	}
	if err := h.lessonTypeAssignmentService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear lesson type assignments: %w", err))
		return
	}
	if err := h.academicRankService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear academic ranks: %w", err))
		return
	}
	if err := h.teacherService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear teachers: %w", err))
		return
	}
	if err := h.disciplineService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear disciplines: %w", err))
		return
	}
	if err := h.lessonTypeService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear lesson types: %w", err))
		return
	}
	if err := h.studentService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear students: %w", err))
		return
	}
	if err := h.studentGroupService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear student groups: %w", err))
		return
	}
	if err := h.groupCohortService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear group cohorts: %w", err))
		return
	}
	if err := h.groupCohortAssignmentService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear group cohort assignments: %w", err))
		return
	}
	if err := h.studyLoadService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear study loads: %w", err))
		return
	}
	if err := h.lessonSlotService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear lesson slots: %w", err))
		return
	}
	if err := h.lessonOccurrenceService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear lesson occurrences: %w", err))
		return
	}
	if err := h.semesterService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear semesters: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
func (h *databaseHandler) ExtractDataFromGenerator(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var studyLoads []entities.StudyLoad
	url := h.scheduleGeneratorDomain + "/get-study-loads"

	if err := h.doGetAndDecode(ctx, url, r, &studyLoads); err != nil {
		jsonutil.ResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.studyLoadService.AddMultiple(ctx, studyLoads); err != nil {
		jsonutil.ResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	var externalL []entities.ExternalLesson
	url = h.scheduleGeneratorDomain + "/get-lessons"

	if err := h.doGetAndDecode(ctx, url, r, &externalL); err != nil {
		jsonutil.ResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.lessonOccurrenceService.AddFromExternal(ctx, externalL); err != nil {
		jsonutil.ResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	jsonutil.ResponseWithJSON(w, http.StatusOK, nil)
}
func (h *databaseHandler) LoadDataIntoGenerator(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	semesterIDs := []uuid.UUID{}

	if err := json.NewDecoder(r.Body).Decode(&semesterIDs); err != nil {
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"LoadDataIntoGenerator",
			fmt.Errorf("failed to decode semester uuids: %w", err),
			logger.HandlerBadRequest,
		))
		return
	}

	disciplines, err := h.semesterDisciplineService.GetBySemesterIDs(ctx, semesterIDs)
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"LoadDataIntoGenerator",
			fmt.Errorf("failed to form disciplines: %w", err),
			logger.HandlerInternalError,
		))
		return
	}

	lessonTypeAssignments, err := h.lessonTypeAssignmentService.GetByDisciplineIDs(
		ctx, h.disciplineService.ExtractIDs(disciplines))
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"LoadDataIntoGenerator",
			fmt.Errorf("failed to form lesson type assignments: %w", err),
			logger.HandlerInternalError,
		))
		return
	}

	lessonTypes, err := h.lessonTypeService.GetMultipleByIDs(
		ctx, h.lessonTypeAssignmentService.GetUniqueLessonTypeIDs(lessonTypeAssignments))
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"LoadDataIntoGenerator",
			fmt.Errorf("failed to form lesson types: %w", err),
			logger.HandlerInternalError,
		))
		return
	}

	groupCohorts, err := h.groupCohortService.GetBySemesterIDs(ctx, semesterIDs)
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"LoadDataIntoGenerator",
			fmt.Errorf("failed to form group cohorts: %w", err),
			logger.HandlerInternalError,
		))
		return
	}

	groupCohortAssignments, err := h.groupCohortAssignmentService.GetByGroupCohortIDs(
		ctx, h.groupCohortService.GetUniqueGroupCohortIDs(groupCohorts))
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"LoadDataIntoGenerator",
			fmt.Errorf("failed to form group cohort assignments: %w", err),
			logger.HandlerInternalError,
		))
		return
	}

	teacherLoads, err := h.teacherLoadService.GetByLessonTypeAssignments(ctx, lessonTypeAssignments)
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"LoadDataIntoGenerator",
			fmt.Errorf("failed to form teacher loads: %w", err),
			logger.HandlerInternalError,
		))
		return
	}

	teachers, err := h.teacherService.GetFullTeachersByIDs(ctx, h.teacherLoadService.GetUniqueTeacherIDs(teacherLoads))
	if err != nil {
		jsonutil.ResponseWithError(w, 500, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"LoadDataIntoGenerator",
			fmt.Errorf("failed to form teachers: %w", err),
			logger.HandlerInternalError,
		))
		return
	}

	url := h.scheduleGeneratorDomain + "/set-disciplines"
	resp := map[string]any{}

	generatorDisciplines := h.disciplineService.ToGeneratorDisciplines(ctx, disciplines)
	if status, err := h.doPostAndDecode(ctx, url, r, generatorDisciplines, &resp); err != nil {
		if status == 400 {
			err = fmt.Errorf(resp["error"].(string))
		}
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"LoadDataIntoGenerator",
			fmt.Errorf("failed to load disciplines to generator: %w", err),
			logger.HandlerInternalError,
		))
		return
	}

	url = h.scheduleGeneratorDomain + "/set-lesson-types"
	resp = map[string]any{}

	generatorLessonTypes := h.lessonTypeService.ToGeneratorLessonType(ctx, lessonTypes)
	if status, err := h.doPostAndDecode(ctx, url, r, generatorLessonTypes, &resp); err != nil {
		if status == 400 {
			err = fmt.Errorf(resp["error"].(string))
		}
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"LoadDataIntoGenerator",
			fmt.Errorf("failed to load lesson types to generator: %w", err),
			logger.HandlerInternalError,
		))
		return
	}

	url = h.scheduleGeneratorDomain + "/set-lesson-type-assignments"
	resp = map[string]any{}

	generatorLessonTypeAssignments := h.lessonTypeAssignmentService.ToGeneratorLessonTypeAssignments(ctx, lessonTypeAssignments)
	if status, err := h.doPostAndDecode(ctx, url, r, generatorLessonTypeAssignments, &resp); err != nil {
		if status == 400 {
			err = fmt.Errorf(resp["error"].(string))
		}
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"LoadDataIntoGenerator",
			fmt.Errorf("failed to load lesson type assignments to generator: %w", err),
			logger.HandlerInternalError,
		))
		return
	}

	url = h.scheduleGeneratorDomain + "/set-student-groups"
	resp = map[string]any{}
	body := struct {
		GroupCohorts           []services.GeneratorGroupCohort           `json:"group_cohorts"`
		GroupCohortAssignments []services.GeneratorGroupCohortAssignment `json:"group_cohort_assignments"`
	}{}
	body.GroupCohorts = h.groupCohortService.ToGeneratorGroupCohorts(ctx, groupCohorts)
	body.GroupCohortAssignments = h.groupCohortAssignmentService.ToGeneratorGroupCohortAssignments(ctx, groupCohortAssignments)
	if status, err := h.doPostAndDecode(ctx, url, r, body, &resp); err != nil {
		if status == 400 {
			err = fmt.Errorf(resp["error"].(string))
		}
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"LoadDataIntoGenerator",
			fmt.Errorf("failed to load student groups to generator: %w", err),
			logger.HandlerInternalError,
		))
		return
	}

	url = h.scheduleGeneratorDomain + "/set-teachers"
	resp = map[string]any{}

	generatorTeachers := h.teacherService.ToGeneratorTeachers(ctx, teachers)
	if status, err := h.doPostAndDecode(ctx, url, r, generatorTeachers, &resp); err != nil {
		if status == 400 {
			err = fmt.Errorf(resp["error"].(string))
		}
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"LoadDataIntoGenerator",
			fmt.Errorf("failed to load lesson type assignments to generator: %w", err),
			logger.HandlerInternalError,
		))
		return
	}

	url = h.scheduleGeneratorDomain + "/set-study-loads"
	resp = map[string]any{}

	generatorTeacherLoads := h.teacherLoadService.ToGeneratorTeacherLoads(ctx, teacherLoads)
	if status, err := h.doPostAndDecode(ctx, url, r, generatorTeacherLoads, &resp); err != nil {
		if status == 400 {
			err = fmt.Errorf(resp["error"].(string))
		}
		jsonutil.ResponseWithError(w, 400, h.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"LoadDataIntoGenerator",
			fmt.Errorf("failed to load teacher loads to generator: %w", err),
			logger.HandlerInternalError,
		))
		return
	}

	jsonutil.ResponseWithJSON(w, http.StatusNoContent, nil)
}

func (h *databaseHandler) doGetAndDecode(ctx context.Context, url string, r *http.Request, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request %q: %w", url, err)
	}

	for k, vv := range r.Header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
func (h *databaseHandler) doPostAndDecode(
	ctx context.Context,
	url string,
	r *http.Request,
	body any,
	target any,
) (int, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return 0, fmt.Errorf("failed to create request %q: %w", url, err)
	}

	req.Header.Set("Content-Type", "application/json")

	for k, vv := range r.Header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 400 {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return resp.StatusCode, fmt.Errorf("failed to decode response: %w", err)
	}

	return resp.StatusCode, nil
}
