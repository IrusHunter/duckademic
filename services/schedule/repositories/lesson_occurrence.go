package repositories

import (
	"context"
	"fmt"
	"hash/fnv"
	"sort"
	"time"

	"github.com/IrusHunter/duckademic/services/schedule/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LessonOccurrenceRepository interface {
	platform.BaseRepository[entities.LessonOccurrence]
	GetLessonsForTeacher(ctx context.Context, tID uuid.UUID, startTime, endTime time.Time) ([]entities.LessonOccurrence, error)
	GetLessonsForStudentGroups(
		ctx context.Context, sgIDs []uuid.UUID, startTime, endTime time.Time) ([]entities.LessonOccurrence, error)
	GetAllLessonsForTeacher(ctx context.Context, teacherID uuid.UUID) ([]entities.LessonOccurrence, error)
	GetAllLessonsForStudentGroups(ctx context.Context, sgIDs []uuid.UUID) ([]entities.LessonOccurrence, error)
	GetLessonsCountForGroups(context.Context, *sqlx.Tx, []uuid.UUID, time.Time) (int, error)

	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	CommitTx(tx *sqlx.Tx) error
	RollbackTx(tx *sqlx.Tx) error

	LockTeacherDate(context.Context, *sqlx.Tx, uuid.UUID, time.Time) error
	LockStudentGroupsDate(context.Context, *sqlx.Tx, []uuid.UUID, time.Time) error
	LockClassroomDate(context.Context, *sqlx.Tx, uuid.UUID, time.Time) error
}

func NewLessonOccurrenceRepository(db *sqlx.DB) LessonOccurrenceRepository {
	config := platform.NewRepositoryConfig(
		"LessonOccurrenceRepository",
		entities.LessonOccurrence{}.TableName(),
		entities.LessonOccurrence{}.EntityName(),
		[]string{
			"id",
			"study_load_id",
			"teacher_id",
			"student_group_id",
			"lesson_slot_id",
			"date",
			"classroom_id",
			"status",
		},
		[]string{"date", "classroom_id"},
		[]string{"created_at", "updated_at"},
	)

	return &lessonOccurrenceRepository{
		BaseRepository: platform.NewBaseRepository[entities.LessonOccurrence](config, db),
		db:             db,
	}
}

type lessonOccurrenceRepository struct {
	platform.BaseRepository[entities.LessonOccurrence]
	db *sqlx.DB
}

type LessonOccurrenceFlat struct {
	ID           uuid.UUID  `db:"id"`
	Date         time.Time  `db:"date"`
	Status       string     `db:"status"`
	MovedToID    *uuid.UUID `db:"moved_to_id"`
	MovedFromID  *uuid.UUID `db:"moved_from_id"`
	StudyLoadID  uuid.UUID  `db:"study_load_id"`
	LessonSlotID uuid.UUID  `db:"lesson_slot_id"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`

	TeacherID        uuid.UUID `db:"teacher_id"`
	TeacherName      string    `db:"teacher_name"`
	StudentGroupID   uuid.UUID `db:"student_group_id"`
	StudentGroupName string    `db:"student_group_name"`
	DisciplineID     uuid.UUID `db:"discipline_id"`
	DisciplineName   string    `db:"discipline_name"`
	LessonTypeID     uuid.UUID `db:"lesson_type_id"`
	LessonTypeName   string    `db:"lesson_type_name"`

	ClassroomID        *uuid.UUID `db:"classroom_id"`
	ClassroomSlug      string     `db:"classroom_slug"`
	ClassroomNumber    string     `db:"classroom_number"`
	ClassroomCapacity  int        `db:"classroom_capacity"`
	ClassroomCreatedAt time.Time  `db:"classroom_created_at"`
	ClassroomUpdatedAt time.Time  `db:"classroom_updated_at"`
}

func (f *LessonOccurrenceFlat) Convert() entities.LessonOccurrence {
	var classroom *entities.Classroom

	if f.ClassroomID != nil {
		classroom = &entities.Classroom{
			ID:        *f.ClassroomID,
			Slug:      f.ClassroomSlug,
			Number:    f.ClassroomNumber,
			Capacity:  f.ClassroomCapacity,
			CreatedAt: f.ClassroomCreatedAt,
			UpdatedAt: f.ClassroomUpdatedAt,
		}
	}

	return entities.LessonOccurrence{
		ID:           f.ID,
		Date:         f.Date,
		ClassroomID:  f.ClassroomID,
		Status:       entities.LessonOccurrenceStatus(f.Status),
		MovedToID:    f.MovedToID,
		StudyLoadID:  f.StudyLoadID,
		LessonSlotID: f.LessonSlotID,
		MovedFromID:  f.MovedFromID,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
		StudyLoad: &entities.StudyLoad{
			ID:               f.StudyLoadID,
			TeacherID:        f.TeacherID,
			TeacherName:      f.TeacherName,
			StudentGroupID:   f.StudentGroupID,
			StudentGroupName: f.StudentGroupName,
			DisciplineID:     f.DisciplineID,
			DisciplineName:   f.DisciplineName,
			LessonTypeID:     f.LessonTypeID,
			LessonTypeName:   f.LessonTypeName,
		},
		Classroom: classroom,
	}
}

func (r *lessonOccurrenceRepository) GetLessonsForTeacher(
	ctx context.Context,
	teacherID uuid.UUID,
	startTime, endTime time.Time,
) ([]entities.LessonOccurrence, error) {
	query := fmt.Sprintf(`
		SELECT 
			lo.id,
			lo.date,
			lo.classroom_id,
			lo.status,
			lo.moved_to_id,
			lo.study_load_id,
			lo.lesson_slot_id,
			lo.moved_from_id,
			lo.created_at,
			lo.updated_at,

			sl.teacher_id,
			sl.teacher_name,
			sl.student_group_id,
			sl.student_group_name,
			sl.discipline_id,
			sl.discipline_name,
			sl.lesson_type_id,
			sl.lesson_type_name,

			c.id as classroom_id,
			c.slug as classroom_slug,
			c.number as classroom_number,
			c.capacity as classroom_capacity,
			c.created_at as classroom_created_at,
			c.updated_at as classroom_updated_at

		FROM %s lo
		JOIN %s sl ON lo.study_load_id = sl.id
		LEFT JOIN %s c ON lo.classroom_id = c.id

		WHERE lo.teacher_id = $1
		  AND lo.date BETWEEN $2 AND $3
		ORDER BY lo.date;
	`, entities.LessonOccurrence{}.TableName(), entities.StudyLoad{}.TableName(), entities.Classroom{}.TableName())

	var flats []LessonOccurrenceFlat
	if err := r.db.SelectContext(ctx, &flats, query, teacherID, startTime, endTime); err != nil {
		return nil, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetLessonsForTeacher",
			err,
			logger.RepositoryScanFailed,
		)
	}

	result := make([]entities.LessonOccurrence, 0, len(flats))
	for i := range flats {
		result = append(result, flats[i].Convert())
	}

	return result, nil
}
func (r *lessonOccurrenceRepository) GetLessonsForStudentGroups(
	ctx context.Context,
	studentGroupIDs []uuid.UUID,
	startTime, endTime time.Time,
) ([]entities.LessonOccurrence, error) {
	if len(studentGroupIDs) == 0 {
		return []entities.LessonOccurrence{}, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetLessonsForStudentGroups",
			fmt.Errorf("given slice of student group ids is empty"),
			logger.RepositoryQueryFailed,
		)
	}

	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT 
			lo.id,
			lo.date,
			lo.classroom_id,
			lo.status,
			lo.moved_to_id,
			lo.study_load_id,
			lo.lesson_slot_id,
			lo.moved_from_id,
			lo.created_at,
			lo.updated_at,

			sl.teacher_id,
			sl.teacher_name,
			sl.student_group_id,
			sl.student_group_name,
			sl.discipline_id,
			sl.discipline_name,
			sl.lesson_type_id,
			sl.lesson_type_name,

			c.id as classroom_id,
			c.slug as classroom_slug,
			c.number as classroom_number,
			c.capacity as classroom_capacity,
			c.created_at as classroom_created_at,
			c.updated_at as classroom_updated_at

		FROM %s lo
		JOIN %s sl ON lo.study_load_id = sl.id
		LEFT JOIN %s c ON lo.classroom_id = c.id

		WHERE sl.student_group_id IN (?)
		  AND lo.date BETWEEN ? AND ?
		ORDER BY lo.date;
	`, entities.LessonOccurrence{}.TableName(), entities.StudyLoad{}.TableName(), entities.Classroom{}.TableName()),
		studentGroupIDs,
		startTime,
		endTime,
	)
	if err != nil {
		return nil, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetLessonsForStudentGroups",
			err,
			logger.RepositoryQueryFailed,
		)
	}

	query = r.db.Rebind(query)

	var flats []LessonOccurrenceFlat
	if err := r.db.SelectContext(ctx, &flats, query, args...); err != nil {
		return nil, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetLessonsForStudentGroups",
			err,
			logger.RepositoryScanFailed,
		)
	}

	result := make([]entities.LessonOccurrence, 0, len(flats))
	for i := range flats {
		result = append(result, flats[i].Convert())
	}

	return result, nil
}
func (r *lessonOccurrenceRepository) GetAllLessonsForTeacher(
	ctx context.Context,
	teacherID uuid.UUID,
) ([]entities.LessonOccurrence, error) {
	query := fmt.Sprintf(`
		SELECT
			lo.id,
			lo.date,
			lo.classroom_id,
			lo.status,
			lo.moved_to_id,
			lo.study_load_id,
			lo.lesson_slot_id,
			lo.moved_from_id,
			lo.created_at,
			lo.updated_at,

			sl.teacher_id,
			sl.teacher_name,
			sl.student_group_id,
			sl.student_group_name,
			sl.discipline_id,
			sl.discipline_name,
			sl.lesson_type_id,
			sl.lesson_type_name,

			c.id as classroom_id,
			c.slug as classroom_slug,
			c.number as classroom_number,
			c.capacity as classroom_capacity,
			c.created_at as classroom_created_at,
			c.updated_at as classroom_updated_at

		FROM %s lo
		JOIN %s sl ON lo.study_load_id = sl.id
		LEFT JOIN %s c ON lo.classroom_id = c.id

		WHERE lo.teacher_id = $1
		ORDER BY lo.date;
	`, entities.LessonOccurrence{}.TableName(), entities.StudyLoad{}.TableName(), entities.Classroom{}.TableName())

	var flats []LessonOccurrenceFlat
	if err := r.db.SelectContext(ctx, &flats, query, teacherID); err != nil {
		return nil, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetAllLessonsForTeacher",
			err,
			logger.RepositoryScanFailed,
		)
	}

	result := make([]entities.LessonOccurrence, 0, len(flats))
	for i := range flats {
		result = append(result, flats[i].Convert())
	}

	return result, nil
}

func (r *lessonOccurrenceRepository) GetAllLessonsForStudentGroups(
	ctx context.Context,
	studentGroupIDs []uuid.UUID,
) ([]entities.LessonOccurrence, error) {
	if len(studentGroupIDs) == 0 {
		return []entities.LessonOccurrence{}, nil
	}

	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT
			lo.id,
			lo.date,
			lo.classroom_id,
			lo.status,
			lo.moved_to_id,
			lo.study_load_id,
			lo.lesson_slot_id,
			lo.moved_from_id,
			lo.created_at,
			lo.updated_at,

			sl.teacher_id,
			sl.teacher_name,
			sl.student_group_id,
			sl.student_group_name,
			sl.discipline_id,
			sl.discipline_name,
			sl.lesson_type_id,
			sl.lesson_type_name,

			c.id as classroom_id,
			c.slug as classroom_slug,
			c.number as classroom_number,
			c.capacity as classroom_capacity,
			c.created_at as classroom_created_at,
			c.updated_at as classroom_updated_at

		FROM %s lo
		JOIN %s sl ON lo.study_load_id = sl.id
		LEFT JOIN %s c ON lo.classroom_id = c.id

		WHERE sl.student_group_id IN (?)
		ORDER BY lo.date;
	`, entities.LessonOccurrence{}.TableName(), entities.StudyLoad{}.TableName(), entities.Classroom{}.TableName()),
		studentGroupIDs,
	)
	if err != nil {
		return nil, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetAllLessonsForStudentGroups",
			err,
			logger.RepositoryQueryFailed,
		)
	}

	query = r.db.Rebind(query)

	var flats []LessonOccurrenceFlat
	if err := r.db.SelectContext(ctx, &flats, query, args...); err != nil {
		return nil, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetAllLessonsForStudentGroups",
			err,
			logger.RepositoryScanFailed,
		)
	}

	result := make([]entities.LessonOccurrence, 0, len(flats))
	for i := range flats {
		result = append(result, flats[i].Convert())
	}

	return result, nil
}

func (r *lessonOccurrenceRepository) GetLessonsCountForGroups(
	ctx context.Context,
	tx *sqlx.Tx,
	groupIDs []uuid.UUID,
	date time.Time,
) (int, error) {
	year, month, day := date.Date()
	startDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, month, day, 23, 59, 59, 0, time.UTC)

	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT
			COUNT(lo.id) as lessons_count

		FROM %s lo
		JOIN %s sl
			ON lo.study_load_id = sl.id

		WHERE sl.student_group_id IN (?)
		AND lo.date BETWEEN ? AND ?
		AND lo.status != 'canceled'
	`,
		entities.LessonOccurrence{}.TableName(),
		entities.StudyLoad{}.TableName(),
	),
		groupIDs,
		startDate,
		endDate,
	)

	if err != nil {
		return 0, err
	}

	query = tx.Rebind(query)

	type row struct {
		LessonsCount int `db:"lessons_count"`
	}

	var rows row

	err = tx.GetContext(ctx, &rows, query, args...)
	if err != nil {
		return 0, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetLessonsCountForGroups",
			err,
			logger.RepositoryScanFailed,
		)
	}

	return rows.LessonsCount, nil
}

func (r *lessonOccurrenceRepository) LockTeacherDate(ctx context.Context, tx *sqlx.Tx, tID uuid.UUID, date time.Time) error {
	key := buildLockKey(
		"teacher",
		tID.String(),
		date.Format(db.TimeFormat),
	)

	_, err := tx.ExecContext(
		ctx,
		`SELECT pg_advisory_xact_lock($1)`,
		key,
	)

	if err != nil {
		return r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"LockTeacherSlot",
			err,
			logger.RepositoryQueryFailed,
		)
	}

	return nil
}
func (r *lessonOccurrenceRepository) LockStudentGroupsDate(
	ctx context.Context,
	tx *sqlx.Tx,
	groupIDs []uuid.UUID,
	date time.Time,
) error {
	sort.Slice(groupIDs, func(i, j int) bool {
		return groupIDs[i].String() < groupIDs[j].String()
	})

	for _, groupID := range groupIDs {
		key := buildLockKey(
			"group",
			groupID.String(),
			date.Format(db.TimeFormat),
		)

		_, err := tx.ExecContext(
			ctx,
			`SELECT pg_advisory_xact_lock($1)`,
			key,
		)

		if err != nil {
			return r.GetLogger().LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"LockStudentGroupsSlot",
				err,
				logger.RepositoryQueryFailed,
			)
		}
	}

	return nil
}
func (r *lessonOccurrenceRepository) LockClassroomDate(
	ctx context.Context,
	tx *sqlx.Tx,
	classroomID uuid.UUID,
	date time.Time,
) error {
	key := buildLockKey(
		"classroom",
		classroomID.String(),
		date.Format(db.TimeFormat),
	)

	_, err := tx.ExecContext(
		ctx,
		`SELECT pg_advisory_xact_lock($1)`,
		key,
	)

	if err != nil {
		return r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"LockClassroomSlot",
			err,
			logger.RepositoryQueryFailed,
		)
	}

	return nil
}

func (r *lessonOccurrenceRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"BeginTx",
			err,
			logger.RepositoryQueryFailed,
		)
	}

	return tx, nil
}
func (r *lessonOccurrenceRepository) CommitTx(tx *sqlx.Tx) error {
	return tx.Commit()
}
func (r *lessonOccurrenceRepository) RollbackTx(tx *sqlx.Tx) error {
	return tx.Rollback()
}

func buildLockKey(parts ...string) int64 {
	h := fnv.New64a()

	for _, p := range parts {
		_, _ = h.Write([]byte(p))
	}

	return int64(h.Sum64())
}
