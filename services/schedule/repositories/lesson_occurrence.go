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
		ctx context.Context,
		teacherID uuid.UUID,
		startTime,
		endTime time.Time,
	) []entities.ScheduledLesson
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

type ScheduledLessonFlat struct {
	ID          uuid.UUID  `db:"id"`
	Date        time.Time  `db:"date"`
	ClassroomID *uuid.UUID `db:"classroom_id"`
	Status      string     `db:"status"`
	MovedToID   *uuid.UUID `db:"moved_to_id"`

	TeacherID        uuid.UUID `db:"teacher_id"`
	TeacherName      string    `db:"teacher_name"`
	StudentGroupID   uuid.UUID `db:"student_group_id"`
	StudentGroupName string    `db:"student_group_name"`
	DisciplineID     uuid.UUID `db:"discipline_id"`
	DisciplineName   string    `db:"discipline_name"`
	LessonTypeID     uuid.UUID `db:"lesson_type_id"`
	LessonTypeName   string    `db:"lesson_type_name"`
}

func (f *ScheduledLessonFlat) Convert() entities.ScheduledLesson {
	return entities.ScheduledLesson{
		ID:          f.ID,
		Date:        f.Date,
		ClassroomID: f.ClassroomID,
		Status:      entities.LessonOccurrenceStatus(f.Status),
		MovedToID:   f.MovedToID,
		StudyLoad: entities.CompactStudyLoad{
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
) []entities.ScheduledLesson {

	query := fmt.Sprintf(`
		SELECT 
			lo.id,
			lo.date,
			lo.classroom_id,
			lo.status,
			lo.moved_to_id,

			sl.teacher_id,
			t.name AS teacher_name,

			sl.student_group_id,
			sg.name AS student_group_name,

			sl.discipline_id,
			d.name AS discipline_name,

			sl.lesson_type_id,
			lt.name AS lesson_type_name

		FROM %s lo
		JOIN %s sl ON lo.study_load_id = sl.id
		LEFT JOIN %s t ON sl.teacher_id = t.id
		LEFT JOIN %s sg ON sl.student_group_id = sg.id
		LEFT JOIN %s d ON sl.discipline_id = d.id
		LEFT JOIN %s lt ON sl.lesson_type_id = lt.id

		WHERE sl.teacher_id = $1
		  AND lo.date BETWEEN $2 AND $3
		ORDER BY lo.date;
	`, entities.LessonOccurrence{}.TableName(), entities.StudyLoad{}.TableName(), entities.Teacher{}.TableName(),
		entities.StudentGroup{}.TableName(), entities.Discipline{}.TableName(), entities.LessonType{}.TableName())

	var flats []ScheduledLessonFlat
	if err := r.db.SelectContext(ctx, &flats, query, teacherID, startTime, endTime); err != nil {
		r.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"GetLessonsForTeacher",
			err,
			logger.RepositoryScanFailed,
		)
		return nil
	}

	result := make([]entities.ScheduledLesson, 0, len(flats))
	for i := range flats {
		result = append(result, flats[i].Convert())
	}

	return result
}
