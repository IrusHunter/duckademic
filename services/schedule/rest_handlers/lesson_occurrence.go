package resthandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
)

// LessonOccurrenceHandler represents a handler responsible for LessonOccurrence-related HTTP operations.
type LessonOccurrenceHandler interface {
	platform.BaseHandler[entities.LessonOccurrence]
	GetPersonalSchedule(context.Context, http.ResponseWriter, *http.Request)
}

// NewLessonOccurrenceHandler creates a new LessonOccurrenceHandler instance.
//
// It requires a lesson occurrence service.
func NewLessonOccurrenceHandler(los services.LessonOccurrenceService) LessonOccurrenceHandler {
	hc := platform.NewHandlerConfig("LessonOccurrenceHandler", entities.LessonOccurrence{}.EntityName())

	return &lessonOccurrenceHandler{
		BaseHandler: platform.NewBaseHandler(hc, los),
		service:     los,
	}
}

type lessonOccurrenceHandler struct {
	platform.BaseHandler[entities.LessonOccurrence]
	service services.LessonOccurrenceService
}

func (h *lessonOccurrenceHandler) GetPersonalSchedule(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	claims := contextutil.GetAccessClaims(ctx)
	if claims == nil {
		jsonutil.ResponseWithError(w, http.StatusUnauthorized, h.GetLogger().LogAndReturnError(contextutil.GetTraceID(ctx),
			"GetPersonalSchedule", fmt.Errorf("failed to get user claims"), logger.HandlerBadRequest))
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		jsonutil.ResponseWithError(w, http.StatusUnauthorized, h.GetLogger().LogAndReturnError(contextutil.GetTraceID(ctx),
			"GetPersonalSchedule", fmt.Errorf("failed to parse user id: %w", err), logger.HandlerBadRequest))
		return
	}

	body := struct {
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		jsonutil.ResponseWithError(w, http.StatusBadRequest, h.GetLogger().LogAndReturnError(contextutil.GetTraceID(ctx),
			"GetPersonalSchedule", fmt.Errorf("failed to extract body: %w", err), logger.HandlerBadRequest))
		return
	}
	defer r.Body.Close()

	var lessons []entities.LessonOccurrence
	switch claims.Role {
	case "student":
		lessons, err = h.service.GetLessonsForStudent(ctx, userID, body.StartTime, body.EndTime)
	case "teacher":
		lessons, err = h.service.GetLessonsForTeacher(ctx, userID, body.StartTime, body.EndTime)
	default:
		jsonutil.ResponseWithError(w, http.StatusBadRequest, h.GetLogger().LogAndReturnError(contextutil.GetTraceID(ctx),
			"GetPersonalSchedule", fmt.Errorf("unknown role %q", claims.Role), logger.HandlerBadRequest))
		return
	}

	if err != nil {
		jsonutil.ResponseWithError(w, http.StatusInternalServerError,
			h.GetLogger().LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"GetPersonalSchedule",
				err,
				logger.HandlerInternalError,
			),
		)
		return
	}

	jsonutil.ResponseWithJSON(w, http.StatusOK, lessons)
}
