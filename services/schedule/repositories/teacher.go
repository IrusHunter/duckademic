package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TeacherRepository interface {
	platform.BaseRepository[entities.Teacher]
	FindByName(context.Context, string) *entities.Teacher
	ExternalUpdate(context.Context, uuid.UUID, entities.Teacher) (entities.Teacher, error)
	Fill(context.Context, uuid.UUID) *entities.Teacher
}

func NewTeacherRepository(db *sqlx.DB) TeacherRepository {
	config := platform.NewRepositoryConfig("TeacherRepository", entities.Teacher{}.TableName(),
		"teacher", []string{"id", "slug", "name", "academic_rank_id"}, []string{""}, []string{"created_at", "updated_at"},
	)
	return &teacherRepository{
		BaseRepository: platform.NewBaseRepository[entities.Teacher](config, db),
		db:             db,
	}
}

type teacherRepository struct {
	platform.BaseRepository[entities.Teacher]
	db *sqlx.DB
}

func (r *teacherRepository) FindByName(ctx context.Context, name string) *entities.Teacher {
	return r.FindFirstBy(ctx, "name", name)
}
func (r *teacherRepository) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	teacher entities.Teacher,
) (entities.Teacher, error) {
	return r.UpdateFields(ctx, id, []string{"slug", "name", "academic_rank_id"}, teacher)
}

type TeacherFlat struct {
	ID             uuid.UUID  `db:"id"`
	Slug           string     `db:"slug"`
	Name           string     `db:"name"`
	AcademicRankID uuid.UUID  `db:"academic_rank_id"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at"`

	AcademicRankTitle    string `db:"academic_rank_title"`
	AcademicRankSlug     string `db:"academic_rank_slug"`
	AcademicRankPriority int    `db:"academic_rank_priority"`
}

func (flat *TeacherFlat) ConvertToTeacher() entities.Teacher {
	academicRank := &entities.AcademicRank{
		ID:       flat.AcademicRankID,
		Slug:     flat.AcademicRankSlug,
		Title:    flat.AcademicRankTitle,
		Priority: flat.AcademicRankPriority,
	}

	return entities.Teacher{
		ID:             flat.ID,
		Slug:           flat.Slug,
		Name:           flat.Name,
		AcademicRankID: flat.AcademicRankID,
		CreatedAt:      flat.CreatedAt,
		UpdatedAt:      flat.UpdatedAt,
		DeletedAt:      flat.DeletedAt,
		AcademicRank:   academicRank,
	}
}

func (r *teacherRepository) Fill(ctx context.Context, id uuid.UUID) *entities.Teacher {
	query := fmt.Sprintf(`
		SELECT 
			t.id,
			t.slug,
			t.name,
			t.academic_rank_id,
			t.created_at,
			t.updated_at,
			t.deleted_at,

			ar.id AS academic_rank_id,
			ar.slug AS academic_rank_slug,
			ar.title AS academic_rank_title,
			ar.priority AS academic_rank_priority

		FROM %s t
		LEFT JOIN %s ar ON t.academic_rank_id = ar.id
		WHERE t.id = $1
		LIMIT 1;
	`, entities.Teacher{}.TableName(), entities.AcademicRank{}.TableName())

	var teacher TeacherFlat

	if err := r.db.GetContext(ctx, &teacher, query, id); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			r.GetLogger().Log(
				contextutil.GetTraceID(ctx),
				"Fill",
				fmt.Sprintf("teacher with id %q not found", id),
				logger.RepositoryOperationSuccess,
			)
			return nil
		}

		r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"Fill",
			fmt.Errorf("failed to scan teacher with id %q: %w", id, err),
			logger.RepositoryScanFailed,
		)

		return nil
	}

	result := teacher.ConvertToTeacher()

	r.GetLogger().Log(
		contextutil.GetTraceID(ctx),
		"Fill",
		fmt.Sprintf("teacher %s found with full academic rank", result),
		logger.RepositoryOperationSuccess,
	)

	return &result
}
