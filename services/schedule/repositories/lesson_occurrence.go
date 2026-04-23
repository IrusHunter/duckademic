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

type LessonOccurrenceRepository interface {
	platform.BaseRepository[entities.LessonOccurrence]
	GetLessonsForTeacher(
		ctx context.Context, teacherID uuid.UUID, startTime, endTime time.Time) ([]entities.LessonOccurrence, error)
	GetLessonsForStudentGroups(
		ctx context.Context, sgIDs []uuid.UUID, startTime, endTime time.Time) ([]entities.LessonOccurrence, error)
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
		[]string{},
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
	ClassroomID  *uuid.UUID `db:"classroom_id"`
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
}

func (f *LessonOccurrenceFlat) Convert() entities.LessonOccurrence {
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
			sl.lesson_type_name

		FROM %s lo
		JOIN %s sl ON lo.study_load_id = sl.id

		WHERE sl.teacher_id = $1
		  AND lo.date BETWEEN $2 AND $3
		ORDER BY lo.date;
	`, entities.LessonOccurrence{}.TableName(), entities.StudyLoad{}.TableName())

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
			sl.lesson_type_name

		FROM %s lo
		JOIN %s sl ON lo.study_load_id = sl.id

		WHERE sl.student_group_id IN (?)
		  AND lo.date BETWEEN ? AND ?
		ORDER BY lo.date;
	`, entities.LessonOccurrence{}.TableName(), entities.StudyLoad{}.TableName()),
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
