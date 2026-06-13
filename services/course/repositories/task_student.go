package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/IrusHunter/duckademic/services/course/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TaskStudentRepository interface {
	platform.BaseRepository[entities.TaskStudent]
	GetUpcomingTasksFor(context.Context, uuid.UUID, time.Time, int) ([]entities.Task, error)
	GetTasksForStudentInCourse(ctx context.Context, studentID, courseID uuid.UUID) ([]entities.TaskStudent, error)
	GetByTaskID(ctx context.Context, taskID uuid.UUID) ([]entities.TaskStudent, error)
}

func NewTaskStudentRepository(db *sqlx.DB) TaskStudentRepository {
	config := platform.NewRepositoryConfig(
		"TaskStudentRepository",
		entities.TaskStudent{}.TableName(),
		entities.TaskStudent{}.EntityName(),
		[]string{"id", "task_id", "student_id", "mark", "submission_time", "file_url", "link_url"},
		[]string{"mark", "submission_time", "file_url", "link_url"},
		[]string{"created_at", "updated_at"},
	)

	tsr := &taskStudentRepository{
		BaseRepository: platform.NewBaseRepository[entities.TaskStudent](config, db),
		db:             db,
	}
	tsr.logger = tsr.GetLogger()

	return tsr
}

type taskStudentRepository struct {
	platform.BaseRepository[entities.TaskStudent]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *taskStudentRepository) GetUpcomingTasksFor(
	ctx context.Context,
	studentID uuid.UUID,
	startTime time.Time,
	count int,
) ([]entities.Task, error) {
	query := fmt.Sprintf(`
		SELECT
			t.id,
			t.course_id,
			t.slug,
			t.title,
			t.description,
			t.max_mark,
			t.deadline,
			t.created_at,
			t.updated_at
		FROM %s t
		JOIN %s sc ON t.course_id = sc.course_id
		WHERE sc.student_id = ?
		  AND t.deadline >= ?
		  AND NOT EXISTS (
		    SELECT 1 FROM %s ts
		    WHERE ts.task_id = t.id
		      AND ts.student_id = sc.student_id
		  )
		ORDER BY t.deadline
		LIMIT ?;
	`,
		entities.Task{}.TableName(),
		entities.StudentCourse{}.TableName(),
		entities.TaskStudent{}.TableName(),
	)

	query = r.db.Rebind(query)

	var tasks []entities.Task
	if err := r.db.SelectContext(ctx, &tasks, query, studentID, startTime, count); err != nil {
		return nil, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetUpcomingTasksFor",
			err,
			logger.RepositoryScanFailed,
		)
	}

	return tasks, nil
}

type TaskStudentFlat struct {
	ID             uuid.UUID  `db:"id"`
	TaskID         uuid.UUID  `db:"task_id"`
	StudentID      uuid.UUID  `db:"student_id"`
	Mark           *float64   `db:"mark"`
	SubmissionTime *time.Time `db:"submission_time"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`

	CourseID    uuid.UUID `db:"course_id"`
	Slug        string    `db:"slug"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	MaxMark     float64   `db:"max_mark"`
	Deadline    time.Time `db:"deadline"`
}

func (f *TaskStudentFlat) Convert() entities.TaskStudent {
	return entities.TaskStudent{
		ID:             f.ID,
		TaskID:         f.TaskID,
		StudentID:      f.StudentID,
		Mark:           f.Mark,
		SubmissionTime: f.SubmissionTime,
		CreatedAt:      f.CreatedAt,
		UpdatedAt:      f.UpdatedAt,
		Task: &entities.Task{
			ID:          f.TaskID,
			CourseID:    f.CourseID,
			Slug:        f.Slug,
			Title:       f.Title,
			Description: f.Description,
			MaxMark:     f.MaxMark,
			Deadline:    f.Deadline,
		},
	}
}

type taskStudentWithStudentFlat struct {
	ID             uuid.UUID  `db:"id"`
	TaskID         uuid.UUID  `db:"task_id"`
	StudentID      uuid.UUID  `db:"student_id"`
	Mark           *float64   `db:"mark"`
	SubmissionTime *time.Time `db:"submission_time"`
	FileURL        *string    `db:"file_url"`
	LinkURL        *string    `db:"link_url"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	StudentName    string     `db:"student_name"`
}

func (r *taskStudentRepository) GetByTaskID(ctx context.Context, taskID uuid.UUID) ([]entities.TaskStudent, error) {
	query := fmt.Sprintf(`
		SELECT
			ts.id, ts.task_id, ts.student_id, ts.mark, ts.submission_time,
			ts.file_url, ts.link_url, ts.created_at, ts.updated_at,
			s.name AS student_name
		FROM %s ts
		JOIN %s s ON ts.student_id = s.id
		WHERE ts.task_id = $1
		ORDER BY ts.created_at ASC;
	`, entities.TaskStudent{}.TableName(), entities.Student{}.TableName())

	var flats []taskStudentWithStudentFlat
	if err := r.db.SelectContext(ctx, &flats, query, taskID); err != nil {
		return nil, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx), "GetByTaskID", err, logger.RepositoryScanFailed,
		)
	}

	result := make([]entities.TaskStudent, 0, len(flats))
	for _, f := range flats {
		result = append(result, entities.TaskStudent{
			ID:             f.ID,
			TaskID:         f.TaskID,
			StudentID:      f.StudentID,
			Mark:           f.Mark,
			SubmissionTime: f.SubmissionTime,
			FileURL:        f.FileURL,
			LinkURL:        f.LinkURL,
			CreatedAt:      f.CreatedAt,
			UpdatedAt:      f.UpdatedAt,
			StudentName:    f.StudentName,
		})
	}
	return result, nil
}

func (r *taskStudentRepository) GetTasksForStudentInCourse(
	ctx context.Context, studentID, courseID uuid.UUID,
) ([]entities.TaskStudent, error) {

	query := fmt.Sprintf(`
		SELECT 
			ts.id,
			ts.task_id,
			ts.student_id,
			ts.mark,
			ts.submission_time,
			ts.created_at,
			ts.updated_at,

			t.course_id,
			t.slug,
			t.title,
			t.description,
			t.max_mark,
			t.deadline

		FROM %s ts
		JOIN %s t ON ts.task_id = t.id

		WHERE ts.student_id = $1
		  AND t.course_id = $2
		ORDER BY t.deadline;
	`,
		entities.TaskStudent{}.TableName(),
		entities.Task{}.TableName(),
	)

	var flats []TaskStudentFlat
	if err := r.db.SelectContext(ctx, &flats, query, studentID, courseID); err != nil {
		return nil, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetTasksForStudentInCourse",
			err,
			logger.RepositoryScanFailed,
		)
	}

	result := make([]entities.TaskStudent, 0, len(flats))
	for i := range flats {
		result = append(result, flats[i].Convert())
	}

	return result, nil
}
