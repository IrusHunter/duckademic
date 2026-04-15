package resthandlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
)

// DatabaseHandler represents a handler responsible for database-related HTTP operations.
type DatabaseHandler interface {
	// Performs database seeding operations, initializing required data.
	Seed(context.Context, http.ResponseWriter, *http.Request)
	Clear(context.Context, http.ResponseWriter, *http.Request)
}

func NewDatabaseHandler(
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
) DatabaseHandler {
	return &databaseHandler{
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
	}
}

type databaseHandler struct {
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

	jsonutil.ResponseWithJSON(w, 204, nil)
}
