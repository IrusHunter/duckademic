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

type TeacherLoadRepository interface {
	platform.BaseRepository[entities.TeacherLoad]
	ExternalUpdate(context.Context, uuid.UUID, entities.TeacherLoad) (entities.TeacherLoad, error)
	GetByDisciplineAndLessonTypeID(
		ctx context.Context, disciplineID uuid.UUID, LessonTypeID uuid.UUID) ([]entities.TeacherLoad, error)
}

func NewTeacherLoadRepository(db *sqlx.DB) TeacherLoadRepository {
	config := platform.NewRepositoryConfig(
		"TeacherLoadRepository",
		entities.TeacherLoad{}.TableName(),
		"teacher_load",
		[]string{"id", "teacher_id", "discipline_id", "lesson_type_id", "group_count"},
		[]string{},
		[]string{"created_at", "updated_at"},
	)

	tr := &teacherLoadRepository{
		BaseRepository: platform.NewBaseRepository[entities.TeacherLoad](config, db),
		db:             db,
	}
	tr.logger = tr.GetLogger()

	return tr
}

type teacherLoadRepository struct {
	platform.BaseRepository[entities.TeacherLoad]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *teacherLoadRepository) ExternalUpdate(
	ctx context.Context, id uuid.UUID, load entities.TeacherLoad,
) (entities.TeacherLoad, error) {
	return r.UpdateFields(ctx, id, []string{"teacher_id", "discipline_id", "lesson_type_id", "group_count"}, load)
}
func (r *teacherLoadRepository) GetByDisciplineAndLessonTypeID(
	ctx context.Context,
	disciplineID uuid.UUID,
	lessonTypeID uuid.UUID,
) ([]entities.TeacherLoad, error) {
	query := fmt.Sprintf(`
		SELECT id, teacher_id, discipline_id, lesson_type_id, group_count
		FROM %s
		WHERE discipline_id = $1 AND lesson_type_id = $2;
	`, entities.TeacherLoad{}.TableName())

	var loads []entities.TeacherLoad

	if err := r.db.SelectContext(ctx, &loads, query, disciplineID, lessonTypeID); err != nil {
		return nil, r.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetByTeacherAndDisciplineID",
			err,
			logger.RepositoryScanFailed,
		)
	}

	return loads, nil
}
