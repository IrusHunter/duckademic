package repositories

import (
	"context"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LessonTypeAssignmentRepository interface {
	platform.BaseRepository[entities.LessonTypeAssignment]
	ExternalUpdate(context.Context, uuid.UUID, entities.LessonTypeAssignment) (entities.LessonTypeAssignment, error)
	FindByLessonTypeAndDiscipline(ctx context.Context, lessonTypeID, disciplineID uuid.UUID) *entities.LessonTypeAssignment
}

func NewLessonTypeAssignmentRepository(db *sqlx.DB) LessonTypeAssignmentRepository {
	config := platform.NewRepositoryConfig(
		"LessonTypeAssignmentRepository",
		entities.LessonTypeAssignment{}.TableName(),
		"lesson_type_assignment",
		[]string{"id", "lesson_type_id", "discipline_id", "required_hours"},
		[]string{""},
		[]string{"created_at", "updated_at"},
	)

	ltaRepo := &lessonTypeAssignmentRepository{
		BaseRepository: platform.NewBaseRepository[entities.LessonTypeAssignment](config, db),
		db:             db,
	}
	ltaRepo.logger = ltaRepo.GetLogger()

	return ltaRepo
}

type lessonTypeAssignmentRepository struct {
	platform.BaseRepository[entities.LessonTypeAssignment]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *lessonTypeAssignmentRepository) FindByLessonTypeAndDiscipline(
	ctx context.Context, lessonTypeID, disciplineID uuid.UUID,
) *entities.LessonTypeAssignment {
	conditions := map[string]any{
		"lesson_type_id": lessonTypeID,
		"discipline_id":  disciplineID,
	}
	return r.FindFirstByConditions(ctx, conditions)
}
func (r *lessonTypeAssignmentRepository) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	assignment entities.LessonTypeAssignment,
) (entities.LessonTypeAssignment, error) {
	return r.UpdateFields(ctx, id, []string{"required_hours"}, assignment)
}
