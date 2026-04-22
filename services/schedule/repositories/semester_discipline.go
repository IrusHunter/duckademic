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

type SemesterDisciplineRepository interface {
	platform.BaseRepository[entities.SemesterDiscipline]
	GetDisciplinesBySemesterID(context.Context, uuid.UUID) ([]entities.Discipline, error)
}

func NewSemesterDisciplineRepository(db *sqlx.DB) SemesterDisciplineRepository {
	config := platform.NewRepositoryConfig(
		"SemesterDisciplineRepository",
		entities.SemesterDiscipline{}.TableName(),
		"semester_discipline",
		[]string{"id", "semester_id", "discipline_id"},
		[]string{""},
		[]string{"created_at", "updated_at"},
	)

	sr := &semesterDisciplineRepository{
		BaseRepository: platform.NewBaseRepository[entities.SemesterDiscipline](config, db),
		db:             db,
	}
	sr.logger = sr.GetLogger()

	return sr
}

type semesterDisciplineRepository struct {
	platform.BaseRepository[entities.SemesterDiscipline]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *semesterDisciplineRepository) GetDisciplinesBySemesterID(
	ctx context.Context,
	semesterID uuid.UUID,
) ([]entities.Discipline, error) {
	query := fmt.Sprintf(`
		SELECT d.id, d.slug, d.name
		FROM %s sd
		JOIN %s d ON sd.discipline_id = d.id
		WHERE sd.semester_id = $1;
	`, entities.SemesterDiscipline{}.TableName(), entities.Discipline{}.TableName())

	var disciplines []entities.Discipline

	if err := r.db.SelectContext(ctx, &disciplines, query, semesterID); err != nil {
		return nil, r.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetDisciplinesBySemesterID",
			err,
			logger.RepositoryScanFailed,
		)
	}

	return disciplines, nil
}
