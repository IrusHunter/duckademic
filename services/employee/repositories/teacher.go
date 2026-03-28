package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/services/employee/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TeacherRepository interface {
	platform.BaseRepository[entities.Teacher]
	Fill(context.Context, uuid.UUID) *entities.Teacher
}

func NewTeacherRepository(db *sqlx.DB) TeacherRepository {
	config := platform.NewRepositoryConfig("TeacherRepository", "teachers", "teacher",
		[]string{"employee_id", "email", "academic_degree_id", "academic_rank_id"},
		[]string{"email", "academic_degree_id", "academic_rank_id"},
		[]string{"created_at", "updated_at"},
	)

	ts := &teacherRepository{
		BaseRepository: platform.NewBaseRepository[entities.Teacher](config, db),
		db:             db,
	}
	ts.logger = ts.GetLogger()
	return ts
}

type teacherRepository struct {
	platform.BaseRepository[entities.Teacher]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *teacherRepository) FindByID(ctx context.Context, id uuid.UUID) *entities.Teacher {
	return r.FindFirstBy(ctx, "employee_id", id)
}

type TeacherFlat struct {
	TeacherID        uuid.UUID  `db:"teacher_id"`
	EmployeeID       uuid.UUID  `db:"employee_id"`
	Email            string     `db:"email"`
	AcademicDegreeID *uuid.UUID `db:"academic_degree_id"`
	AcademicRankID   uuid.UUID  `db:"academic_rank_id"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at"`

	EmployeeFirstName  string  `db:"employee_first_name"`
	EmployeeLastName   string  `db:"employee_last_name"`
	EmployeeMiddleName *string `db:"employee_middle_name"`
	EmployeeSlug       string  `db:"employee_slug"`

	AcademicDegreeTitle *string `db:"academic_degree_title"`

	AcademicRankTitle string `db:"academic_rank_title"`
}

func (flat *TeacherFlat) ConvertToTeacher() entities.Teacher {
	var academicDegree *entities.AcademicDegree
	if flat.AcademicDegreeID != nil {
		academicDegree = &entities.AcademicDegree{
			ID:    *flat.AcademicDegreeID,
			Title: *flat.AcademicDegreeTitle,
		}
	}

	employee := &entities.Employee{
		ID:         flat.EmployeeID,
		FirstName:  flat.EmployeeFirstName,
		LastName:   flat.EmployeeLastName,
		MiddleName: flat.EmployeeMiddleName,
		Slug:       flat.EmployeeSlug,
	}

	academicRank := &entities.AcademicRank{
		ID:    flat.AcademicRankID,
		Title: flat.AcademicRankTitle,
	}

	return entities.Teacher{
		EmployeeID:       flat.EmployeeID,
		Email:            flat.Email,
		AcademicDegreeID: flat.AcademicDegreeID,
		AcademicRankID:   flat.AcademicRankID,
		CreatedAt:        flat.CreatedAt,
		UpdatedAt:        flat.UpdatedAt,
		DeletedAt:        flat.DeletedAt,
		Employee:         employee,
		AcademicDegree:   academicDegree,
		AcademicRank:     academicRank,
	}
}

func (r *teacherRepository) Fill(ctx context.Context, id uuid.UUID) *entities.Teacher {
	query := `
		SELECT 
			t.employee_id AS teacher_id,
			t.employee_id,
			t.email,
			t.academic_degree_id,
			t.academic_rank_id,
			t.created_at,
			t.updated_at,
			t.deleted_at,

			e.id AS employee_id,
			e.first_name AS employee_first_name,
			e.last_name AS employee_last_name,
			e.middle_name AS employee_middle_name,
			e.slug AS employee_slug,

			ad.id AS academic_degree_id,
			ad.title AS academic_degree_title,

			ar.id AS academic_rank_id,
			ar.title AS academic_rank_title

		FROM teachers t
		LEFT JOIN employees e ON t.employee_id = e.id
		LEFT JOIN academic_degrees ad ON t.academic_degree_id = ad.id
		LEFT JOIN academic_ranks ar ON t.academic_rank_id = ar.id
		WHERE t.employee_id = $1
		LIMIT 1;
	`

	var teacher TeacherFlat

	if err := r.db.GetContext(ctx, &teacher, query, id); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			r.logger.Log(contextutil.GetTraceID(ctx), "Fill",
				fmt.Sprintf("teacher with id %q not found", id),
				logger.RepositoryOperationSuccess,
			)
			return nil
		}

		r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "FindFirstDetailedBy",
			fmt.Errorf("failed to scan teacher with id %q: %w", id, err),
			logger.RepositoryScanFailed,
		)
		return nil
	}

	trueTeacher := teacher.ConvertToTeacher()
	r.logger.Log(contextutil.GetTraceID(ctx), "FindFirstDetailedBy",
		fmt.Sprintf("teacher %s found with relations", trueTeacher),
		logger.RepositoryOperationSuccess,
	)

	return &trueTeacher
}
