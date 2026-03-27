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
	gcs services.GroupCohortService,
	sgs services.StudentGroupService,
	gms services.GroupMemberService,
) DatabaseHandler {
	return &databaseHandler{
		groupCohortService:  gcs,
		studentGroupService: sgs,
		groupMembersService: gms,
	}
}

type databaseHandler struct {
	groupCohortService  services.GroupCohortService
	studentGroupService services.StudentGroupService
	groupMembersService services.GroupMemberService
}

func (h *databaseHandler) Seed(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(events.ExternalSeedCooldown)
		ctx = contextutil.SetTraceID(context.Background())
		h.groupCohortService.Seed(ctx)
		ctx = contextutil.SetTraceID(context.Background())
		h.studentGroupService.Seed(ctx)

		time.Sleep(events.ExternalSeedCooldown)
		ctx = contextutil.SetTraceID(context.Background())
		h.groupMembersService.Seed(ctx)
	}()

	jsonutil.ResponseWithJSON(w, 204, nil)
}
func (h *databaseHandler) Clear(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.groupMembersService.Clear(ctx); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to clear group members: %w", err))
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

	jsonutil.ResponseWithJSON(w, 204, nil)
}
