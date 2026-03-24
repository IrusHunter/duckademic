package repositories

import (
	"github.com/IrusHunter/duckademic/services/curriculum/entities"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/jmoiron/sqlx"
)

type SemesterDisciplineRepository interface {
	platform.BaseRepository[entities.SemesterDiscipline]
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

	sdRepo := &semesterDisciplineRepository{
		BaseRepository: platform.NewBaseRepository[entities.SemesterDiscipline](config, db),
		db:             db,
	}
	sdRepo.logger = sdRepo.GetLogger()

	return sdRepo
}

type semesterDisciplineRepository struct {
	platform.BaseRepository[entities.SemesterDiscipline]
	db     *sqlx.DB
	logger logger.Logger
}
