package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TeacherSlotPriorityRepository interface {
	platform.BaseRepository[entities.TeacherSlotPriority]
	GetByTeacherID(context.Context, uuid.UUID) ([]entities.TeacherSlotPriority, error)
}

func NewTeacherSlotPriorityRepository(db *sqlx.DB) TeacherSlotPriorityRepository {
	config := platform.NewRepositoryConfig(
		"TeacherSlotPriorityRepository",
		entities.TeacherSlotPriority{}.TableName(),
		entities.TeacherSlotPriority{}.EntityName(),
		[]string{
			"id",
			"teacher_id",
			"time_slot_id",
			"priority",
		},
		[]string{"priority"},
		[]string{"created_at", "updated_at"},
	)

	return &teacherSlotPriorityRepository{
		BaseRepository: platform.NewBaseRepository[entities.TeacherSlotPriority](config, db),
		db:             db,
	}
}

type teacherSlotPriorityRepository struct {
	platform.BaseRepository[entities.TeacherSlotPriority]
	db *sqlx.DB
}

type TeacherSlotPriorityFlat struct {
	ID         uuid.UUID `db:"id"`
	TeacherID  uuid.UUID `db:"teacher_id"`
	TimeSlotID uuid.UUID `db:"time_slot_id"`
	Priority   string    `db:"priority"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`

	TimeSlotWeekday int `db:"time_slot_weekday"`
	TimeSlotSlot    int `db:"time_slot_slot"`
}

func (f *TeacherSlotPriorityFlat) Convert() entities.TeacherSlotPriority {
	return entities.TeacherSlotPriority{
		ID:         f.ID,
		TeacherID:  f.TeacherID,
		TimeSlotID: f.TimeSlotID,
		Priority:   entities.TeacherSlotPriorityValue(f.Priority),
		CreatedAt:  f.CreatedAt,
		UpdatedAt:  f.UpdatedAt,
		TimeSlot: &entities.LessonSlot{
			ID:      f.TimeSlotID,
			Weekday: time.Weekday(f.TimeSlotWeekday),
			Slot:    f.TimeSlotSlot,
		},
	}
}

func (r *teacherSlotPriorityRepository) GetByTeacherID(
	ctx context.Context,
	teacherID uuid.UUID,
) ([]entities.TeacherSlotPriority, error) {
	query := fmt.Sprintf(`
		SELECT
			tsp.id,
			tsp.teacher_id,
			tsp.time_slot_id,
			tsp.priority,
			tsp.created_at,
			tsp.updated_at,

			ls.weekday AS time_slot_weekday,
			ls.slot AS time_slot_slot

		FROM %s tsp
		JOIN %s ls ON tsp.time_slot_id = ls.id

		WHERE tsp.teacher_id = $1
		ORDER BY ls.weekday, ls.slot;
	`,
		entities.TeacherSlotPriority{}.TableName(),
		entities.LessonSlot{}.TableName(),
	)

	var flats []TeacherSlotPriorityFlat

	if err := r.db.SelectContext(ctx, &flats, query, teacherID); err != nil {
		return nil, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetByTeacherID",
			err,
			logger.RepositoryScanFailed,
		)
	}

	result := make([]entities.TeacherSlotPriority, 0, len(flats))

	for i := range flats {
		result = append(result, flats[i].Convert())
	}

	return result, nil
}
