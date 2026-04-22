package main

import (
	"context"
	"fmt"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/jmoiron/sqlx"
)

// Migrate creates the necessary database schema and triggers for the application.
//
// Returns an error if any table creation or trigger setup fails.
func Migrate(database *sqlx.DB) error {
	migrationsF := []func(*sqlx.Tx) error{
		academicRankMigrations,
		teacherMigrations,
		semesterMigrations,
		disciplineMigrations,
		lessonTypeMigrations,
		lessonTypeAssignmentMigrations,
		studentMigrations,
		teacherLoadMigrations,
		groupCohortMigrations,
		studentGroupMigrations,
		groupCohortAssignmentMigrations,
		groupMembersMigrations,
		classroomMigrations,
		studyLoadMigrations,
		lessonSlotMigrations,
		lessonOccurrenceMigrations,
		semesterDisciplineMigrations,
	}

	for _, f := range migrationsF {
		tx, err := database.Beginx()
		if err != nil {
			return err
		}

		if err := f(tx); err != nil {
			_ = tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func academicRankMigrations(tx *sqlx.Tx) error {
	schema := ` 
	CREATE TABLE IF NOT EXISTS academic_ranks (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		title TEXT NOT NULL,
		priority INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create academic ranks relation: %w", err)
	}

	index := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_academic_ranks_slug
	ON academic_ranks (slug);
	`

	if _, err := tx.Exec(index); err != nil {
		return fmt.Errorf("failed to create academic_ranks slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "academic_ranks"); err != nil {
		return fmt.Errorf("failed to create on update trigger for academic_ranks: %w", err)
	}

	return nil
}
func teacherMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS teachers (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		academic_rank_id UUID NOT NULL REFERENCES academic_ranks(id),
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		deleted_at TIMESTAMP WITH TIME ZONE
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create teachers table: %w", err)
	}

	indexName := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_teachers_name
	ON teachers (name);
	`

	if _, err := tx.Exec(indexName); err != nil {
		return fmt.Errorf("failed to create teachers name index: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "teachers"); err != nil {
		return fmt.Errorf("failed to create on update trigger for teachers: %w", err)
	}

	addSlug := `
	ALTER TABLE teachers
	ADD COLUMN IF NOT EXISTS slug TEXT;
	`
	if _, err := tx.Exec(addSlug); err != nil {
		return fmt.Errorf("failed to add slug column: %w", err)
	}

	dropIndex := `
	DROP INDEX IF EXISTS idx_teachers_name;
	`
	if _, err := tx.Exec(dropIndex); err != nil {
		return fmt.Errorf("failed to drop name index: %w", err)
	}

	indexName = `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_teachers_slug
	ON teachers (slug);
	`

	if _, err := tx.Exec(indexName); err != nil {
		return fmt.Errorf("failed to create teachers slug index: %w", err)
	}

	return nil
}
func disciplineMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS disciplines (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create disciplines table: %w", err)
	}

	createSlugIndex := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_disciplines_slug
	ON disciplines (slug);
	`
	if _, err := tx.Exec(createSlugIndex); err != nil {
		return fmt.Errorf("failed to create disciplines slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "disciplines"); err != nil {
		return fmt.Errorf("failed to create on update trigger for disciplines: %w", err)
	}

	return nil
}
func lessonTypeMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS lesson_types (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		hours_value INT NOT NULL DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create lesson_types table: %w", err)
	}

	createSlugIndex := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_lesson_types_slug
	ON lesson_types (slug);
	`
	if _, err := tx.Exec(createSlugIndex); err != nil {
		return fmt.Errorf("failed to create lesson_types slug index for lesson_types: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "lesson_types"); err != nil {
		return fmt.Errorf("failed to create on update trigger for lesson_types: %w", err)
	}

	addReservedWeeksColumn := `
	ALTER TABLE lesson_types
	ADD COLUMN IF NOT EXISTS reserved_weeks TEXT NOT NULL DEFAULT '';
	`

	if _, err := tx.Exec(addReservedWeeksColumn); err != nil {
		return fmt.Errorf("failed to add  reserved_weeks for lesson_types: %w", err)
	}

	return nil
}
func lessonTypeAssignmentMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS lesson_type_assignments (
		id UUID PRIMARY KEY,
		lesson_type_id UUID NOT NULL,
		discipline_id UUID NOT NULL,
		required_hours INT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		CONSTRAINT uq_lta_lesson_discipline UNIQUE (lesson_type_id, discipline_id),
		CONSTRAINT fk_lesson_type FOREIGN KEY (lesson_type_id) REFERENCES lesson_types(id) ON DELETE CASCADE,
		CONSTRAINT fk_discipline FOREIGN KEY (discipline_id) REFERENCES disciplines(id) ON DELETE CASCADE
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create lesson_type_assignments table: %w", err)
	}

	indexLessonDiscipline := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_lta_lesson_type_discipline
	ON lesson_type_assignments (lesson_type_id, discipline_id);
	`
	if _, err := tx.Exec(indexLessonDiscipline); err != nil {
		return fmt.Errorf("failed to create lesson_type_assignments unique index: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "lesson_type_assignments"); err != nil {
		return fmt.Errorf("failed to create on update trigger for lesson_type_assignments: %w", err)
	}

	return nil
}
func studentMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS students (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		deleted_at TIMESTAMP WITH TIME ZONE
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create students table: %w", err)
	}

	indexSlug := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_students_slug
	ON students (slug);
	`
	if _, err := tx.Exec(indexSlug); err != nil {
		return fmt.Errorf("failed to create students slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "students"); err != nil {
		return fmt.Errorf("failed to create on update trigger for students: %w", err)
	}

	return nil
}
func studentGroupMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS student_groups (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create student_groups table: %w", err)
	}

	indexSlug := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_student_groups_slug
	ON student_groups (slug);
	`
	if _, err := tx.Exec(indexSlug); err != nil {
		return fmt.Errorf("failed to create student_groups slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "student_groups"); err != nil {
		return fmt.Errorf("failed to create on update trigger for student_groups: %w", err)
	}

	groupCohortIdAdd := `
	ALTER TABLE student_groups
	ADD COLUMN IF NOT EXISTS group_cohort_id UUID NOT NULL REFERENCES group_cohorts(id) ON DELETE CASCADE;
	`

	if _, err := tx.Exec(groupCohortIdAdd); err != nil {
		return fmt.Errorf("failed to add group_cohort_id column to student_groups: %w", err)
	}

	return nil
}
func groupMembersMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS group_members (
		id UUID PRIMARY KEY,
		student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
		student_group_id UUID REFERENCES student_groups(id) ON DELETE SET NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create group_members table: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "group_members"); err != nil {
		return fmt.Errorf("failed to create on update trigger for group_members: %w", err)
	}

	return nil
}
func teacherLoadMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS teacher_loads (
		id UUID PRIMARY KEY,
		teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
		discipline_id UUID NOT NULL REFERENCES disciplines(id) ON DELETE CASCADE,
		lesson_type_id UUID NOT NULL REFERENCES lesson_types(id) ON DELETE CASCADE,
		group_cohort_id UUID NOT NULL,
		group_count INT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create teacher_loads table: %w", err)
	}

	createUniqueIndex := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_teacher_loads_unique
	ON teacher_loads (teacher_id, discipline_id, lesson_type_id);
	`
	if _, err := tx.Exec(createUniqueIndex); err != nil {
		return fmt.Errorf("failed to create teacher_loads unique index: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "teacher_loads"); err != nil {
		return fmt.Errorf("failed to create on update trigger for teacher_loads: %w", err)
	}

	// WARNING: it should be removed if group_cohort_id will be needed
	dropColumn := `
	ALTER TABLE teacher_loads
	DROP COLUMN IF EXISTS group_cohort_id;
	`

	if _, err := tx.Exec(dropColumn); err != nil {
		return fmt.Errorf("failed to drop group_cohort_id column: %w", err)
	}

	return nil
}
func groupCohortMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS group_cohorts (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create group_cohorts table: %w", err)
	}

	indexSlug := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_group_cohorts_slug
	ON group_cohorts (slug);
	`
	if _, err := tx.Exec(indexSlug); err != nil {
		return fmt.Errorf("failed to create group_cohorts slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "group_cohorts"); err != nil {
		return fmt.Errorf("failed to create on update trigger for group_cohorts: %w", err)
	}

	addSemesterColumn := `
	ALTER TABLE group_cohorts
	ADD COLUMN IF NOT EXISTS semester_id UUID NULL;
	`
	if _, err := tx.Exec(addSemesterColumn); err != nil {
		return fmt.Errorf("failed to add semester_id column: %w", err)
	}

	if err := db.EnsureForeignKeyTx(
		context.Background(),
		tx,
		"fk_group_cohorts_semester",
		"group_cohorts",
		"semester_id",
		"semesters",
		"id",
	); err != nil {
		return fmt.Errorf("failed to create semester_id foreign key: %w", err)
	}

	return nil
}
func groupCohortAssignmentMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS group_cohort_assignments (
		id UUID PRIMARY KEY,
		group_cohort_id UUID NOT NULL REFERENCES group_cohorts(id) ON DELETE CASCADE,
		discipline_id UUID NOT NULL REFERENCES disciplines(id) ON DELETE CASCADE,
		lesson_type_id UUID NOT NULL REFERENCES lesson_types(id) ON DELETE CASCADE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create group_cohort_assignments table: %w", err)
	}

	indexUnique := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_group_cohort_assignments_unique
	ON group_cohort_assignments (group_cohort_id, discipline_id, lesson_type_id);
	`
	if _, err := tx.Exec(indexUnique); err != nil {
		return fmt.Errorf("failed to create unique index for group_cohort_assignments: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "group_cohort_assignments"); err != nil {
		return fmt.Errorf("failed to create on update trigger for group_cohort_assignments: %w", err)
	}

	return nil
}
func classroomMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS classrooms (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		number TEXT NOT NULL,
		capacity INT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create classrooms table: %w", err)
	}

	indexSlug := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_classrooms_slug
	ON classrooms (slug);
	`
	if _, err := tx.Exec(indexSlug); err != nil {
		return fmt.Errorf("failed to create classrooms slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "classrooms"); err != nil {
		return fmt.Errorf("failed to create on update trigger for classrooms: %w", err)
	}

	return nil
}
func studyLoadMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS study_loads (
		id UUID PRIMARY KEY,
		teacher_id UUID NOT NULL,
		student_group_id UUID NOT NULL,
		discipline_id UUID NOT NULL,
		lesson_type_id UUID NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create study_loads relation: %w", err)
	}

	index := `
	CREATE UNIQUE INDEX IF NOT EXISTS uq_study_load
	ON study_loads (
		teacher_id,
		student_group_id,
		discipline_id,
		lesson_type_id
	);
	`

	if _, err := tx.Exec(index); err != nil {
		return fmt.Errorf("failed to add unique index to study_loads: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "study_loads"); err != nil {
		return fmt.Errorf("failed to create on update trigger for study_loads: %w", err)
	}

	return nil
}
func lessonSlotMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS lesson_slots (
		id UUID PRIMARY KEY,
		slot INTEGER NOT NULL,
		weekday INTEGER NOT NULL,
		start_time BIGINT NOT NULL,
		duration BIGINT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create lesson_slots relation: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "lesson_slots"); err != nil {
		return fmt.Errorf("failed to create on update trigger for lesson_slots: %w", err)
	}

	return nil
}
func lessonOccurrenceMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS lesson_occurrences (
		id UUID PRIMARY KEY,
		study_load_id UUID NOT NULL,
		teacher_id UUID NOT NULL,
		student_group_id UUID NOT NULL,
		lesson_slot_id UUID NOT NULL,
		date TIMESTAMP WITH TIME ZONE NOT NULL,
		classroom_id UUID NULL,
		status TEXT NOT NULL,
		moved_to_id UUID NULL,
		moved_from_id UUID NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create lesson_occurrences relation: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "lesson_occurrences"); err != nil {
		return fmt.Errorf("failed to create on update trigger for lesson_occurrences: %w", err)
	}

	return nil
}
func semesterMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS semesters (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		curriculum_id UUID NOT NULL,
		number INT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create semesters table: %w", err)
	}

	indexSlug := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_semesters_slug
	ON semesters (slug);
	`
	if _, err := tx.Exec(indexSlug); err != nil {
		return fmt.Errorf("failed to create semesters slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "semesters"); err != nil {
		return fmt.Errorf("failed to create on update trigger for semesters: %w", err)
	}

	return nil
}
func semesterDisciplineMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS semester_discipline (
		id UUID PRIMARY KEY,
		semester_id UUID NOT NULL,
		discipline_id UUID NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		CONSTRAINT uq_semester_discipline UNIQUE (semester_id, discipline_id)
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create semester_discipline table: %w", err)
	}

	indexSemester := `
	CREATE INDEX IF NOT EXISTS idx_semester
	ON semester_discipline (semester_id);
	`
	if _, err := tx.Exec(indexSemester); err != nil {
		return fmt.Errorf("failed to create semester_discipline semester index: %w", err)
	}

	if err := db.EnsureForeignKeyTx(
		context.Background(),
		tx,
		"fk_semester_discipline_semester",
		"semester_discipline",
		"semester_id",
		"semesters",
		"id",
	); err != nil {
		return fmt.Errorf("failed to create semester_id foreign key: %w", err)
	}

	if err := db.EnsureForeignKeyTx(
		context.Background(),
		tx,
		"fk_semester_discipline_discipline",
		"semester_discipline",
		"discipline_id",
		"disciplines",
		"id",
	); err != nil {
		return fmt.Errorf("failed to create discipline_id foreign key: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "semester_discipline"); err != nil {
		return fmt.Errorf("failed to create on update trigger for semester_discipline: %w", err)
	}

	return nil
}
