package resthandlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/IrusHunter/duckademic/services/student_group/services"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
)

type DatabaseHandler interface {
	Seed(context.Context, http.ResponseWriter, *http.Request)
	Clear(context.Context, http.ResponseWriter, *http.Request)
}

func NewDatabaseHandler(
	ss services.StudentService,
	semS services.SemesterService,
	gcs services.GroupCohortService,
	sgs services.StudentGroupService,
	gms services.GroupMemberService,
	ds services.DisciplineService,
	ls services.LessonTypeService,
	gcas services.GroupCohortAssignmentService,
) DatabaseHandler {
	return &databaseHandler{
		studentService:               ss,
		semesterService:              semS,
		groupCohortService:           gcs,
		studentGroupService:          sgs,
		groupMembersService:          gms,
		disciplineService:            ds,
		lessonTypeService:            ls,
		groupCohortAssignmentService: gcas,
	}
}

type databaseHandler struct {
	studentService               services.StudentService
	semesterService              services.SemesterService
	groupCohortService           services.GroupCohortService
	studentGroupService          services.StudentGroupService
	groupMembersService          services.GroupMemberService
	disciplineService            services.DisciplineService
	lessonTypeService            services.LessonTypeService
	groupCohortAssignmentService services.GroupCohortAssignmentService
}

func (h *databaseHandler) Seed(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(events.ExternalSeedCooldown)
		ctx = contextutil.SetTraceID(context.Background())
		h.groupCohortService.Seed(ctx)
		ctx = contextutil.SetTraceID(context.Background())
		h.studentGroupService.Seed(ctx)
		ctx = contextutil.SetTraceID(context.Background())
		h.groupCohortAssignmentService.Seed(ctx)

		time.Sleep(events.ExternalSeedCooldown)
		ctx = contextutil.SetTraceID(context.Background())
		h.groupMembersService.Seed(ctx)
	}()

	jsonutil.ResponseWithJSON(w, 204, nil)
}
func (h *databaseHandler) Clear(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.groupCohortAssignmentService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear group cohort assignments: %w", err))
		return
	}
	if err := h.groupMembersService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear group members: %w", err))
		return
	}
	if err := h.studentGroupService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear student groups: %w", err))
		return
	}
	if err := h.studentService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear students: %w", err))
		return
	}
	if err := h.groupCohortService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear group cohorts: %w", err))
		return
	}
	if err := h.semesterService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear semesters: %w", err))
		return
	}
	if err := h.lessonTypeService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear lesson types: %w", err))
		return
	}
	if err := h.disciplineService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear disciplines: %w", err))
		return
	}

	jsonutil.ResponseWithJSON(w, 204, nil)
}
