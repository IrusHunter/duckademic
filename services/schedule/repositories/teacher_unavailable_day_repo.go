package repositories

import (
	"context"
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TeacherUnavailableDayRepository interface {
	platform.BaseRepository[entities.TeacherUnavailableDay]
	GetByTeacherID(context.Context, uuid.UUID) ([]entities.TeacherUnavailableDay, error)
}

func NewTeacherUnavailableDaysRepository(
	db *sqlx.DB,
) TeacherUnavailableDayRepository {
	config := platform.NewRepositoryConfig(
		"TeacherUnavailableDaysRepository",
		entities.TeacherUnavailableDay{}.TableName(),
		entities.TeacherUnavailableDay{}.EntityName(),
		[]string{"id", "teacher_id", "day"},
		[]string{},
		[]string{"created_at", "updated_at"},
	)

	return &teacherUnavailableDayRepository{
		BaseRepository: platform.NewBaseRepository[entities.TeacherUnavailableDay](
			config,
			db,
		),
		db: db,
	}
}

type teacherUnavailableDayRepository struct {
	platform.BaseRepository[entities.TeacherUnavailableDay]
	db *sqlx.DB
}

func (r *teacherUnavailableDayRepository) GetByTeacherID(
	ctx context.Context,
	teacherID uuid.UUID,
) ([]entities.TeacherUnavailableDay, error) {
	query := fmt.Sprintf(`
		SELECT
			id,
			teacher_id,
			day,
			created_at,
			updated_at
		FROM %s
		WHERE teacher_id = $1
		ORDER BY day;
	`,
		entities.TeacherUnavailableDay{}.TableName(),
	)

	var result []entities.TeacherUnavailableDay

	if err := r.db.SelectContext(ctx, &result, query, teacherID); err != nil {
		return nil, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetByTeacherID",
			err,
			logger.RepositoryScanFailed,
		)
	}

	return result, nil
}
