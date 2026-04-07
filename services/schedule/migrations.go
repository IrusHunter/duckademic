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
	migrationsF := []func(*sqlx.DB) error{
		academicRankMigrations,
		teacherMigrations,
		disciplineMigrations,
		lessonTypeMigrations,
		lessonTypeAssignmentMigrations,
		studentMigrations,
		groupMembersMigrations,
		teacherLoadMigrations,
		groupCohortMigrations,
		studentGroupMigrations,
		groupCohortAssignmentMigrations,
	}

	for _, f := range migrationsF {
		err := f(database)
		if err != nil {
			return err
		}
	}

	return nil
}

func academicRankMigrations(database *sqlx.DB) error {
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

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create academic ranks relation: %w", err)
	}

	index := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_academic_ranks_slug
	ON academic_ranks (slug);
	`

	if _, err := database.Exec(index); err != nil {
		return fmt.Errorf("failed to create academic_ranks slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "academic_ranks"); err != nil {
		return fmt.Errorf("failed to create on update trigger for academic_ranks: %w", err)
	}

	return nil
}
func teacherMigrations(database *sqlx.DB) error {
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

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create teachers table: %w", err)
	}

	indexName := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_teachers_name
	ON teachers (name);
	`

	if _, err := database.Exec(indexName); err != nil {
		return fmt.Errorf("failed to create teachers name index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "teachers"); err != nil {
		return fmt.Errorf("failed to create on update trigger for teachers: %w", err)
	}

	addSlug := `
	ALTER TABLE teachers
	ADD COLUMN IF NOT EXISTS slug TEXT;
	`
	if _, err := database.Exec(addSlug); err != nil {
		return fmt.Errorf("failed to add slug column: %w", err)
	}

	dropIndex := `
	DROP INDEX IF EXISTS idx_teachers_name;
	`
	if _, err := database.Exec(dropIndex); err != nil {
		return fmt.Errorf("failed to drop name index: %w", err)
	}

	indexName = `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_teachers_slug
	ON teachers (slug);
	`

	if _, err := database.Exec(indexName); err != nil {
		return fmt.Errorf("failed to create teachers slug index: %w", err)
	}

	return nil
}
func disciplineMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS disciplines (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create disciplines table: %w", err)
	}

	createSlugIndex := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_disciplines_slug
	ON disciplines (slug);
	`
	if _, err := database.Exec(createSlugIndex); err != nil {
		return fmt.Errorf("failed to create disciplines slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "disciplines"); err != nil {
		return fmt.Errorf("failed to create on update trigger for disciplines: %w", err)
	}

	return nil
}
func lessonTypeMigrations(database *sqlx.DB) error {
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

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create lesson_types table: %w", err)
	}

	createSlugIndex := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_lesson_types_slug
	ON lesson_types (slug);
	`
	if _, err := database.Exec(createSlugIndex); err != nil {
		return fmt.Errorf("failed to create lesson_types slug index for lesson_types: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "lesson_types"); err != nil {
		return fmt.Errorf("failed to create on update trigger for lesson_types: %w", err)
	}

	addReservedWeeksColumn := `
	ALTER TABLE lesson_types
	ADD COLUMN IF NOT EXISTS reserved_weeks TEXT NOT NULL DEFAULT '';
	`

	if _, err := database.Exec(addReservedWeeksColumn); err != nil {
		return fmt.Errorf("failed to add  reserved_weeks for lesson_types: %w", err)
	}

	return nil
}
func lessonTypeAssignmentMigrations(database *sqlx.DB) error {
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

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create lesson_type_assignments table: %w", err)
	}

	indexLessonDiscipline := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_lta_lesson_type_discipline
	ON lesson_type_assignments (lesson_type_id, discipline_id);
	`
	if _, err := database.Exec(indexLessonDiscipline); err != nil {
		return fmt.Errorf("failed to create lesson_type_assignments unique index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "lesson_type_assignments"); err != nil {
		return fmt.Errorf("failed to create on update trigger for lesson_type_assignments: %w", err)
	}

	return nil
}
func studentMigrations(database *sqlx.DB) error {
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

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create students table: %w", err)
	}

	indexSlug := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_students_slug
	ON students (slug);
	`
	if _, err := database.Exec(indexSlug); err != nil {
		return fmt.Errorf("failed to create students slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "students"); err != nil {
		return fmt.Errorf("failed to create on update trigger for students: %w", err)
	}

	return nil
}
func studentGroupMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS student_groups (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create student_groups table: %w", err)
	}

	indexSlug := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_student_groups_slug
	ON student_groups (slug);
	`
	if _, err := database.Exec(indexSlug); err != nil {
		return fmt.Errorf("failed to create student_groups slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "student_groups"); err != nil {
		return fmt.Errorf("failed to create on update trigger for student_groups: %w", err)
	}

	groupCohortIdAdd := `
	ALTER TABLE student_groups
	ADD COLUMN IF NOT EXISTS group_cohort_id UUID NOT NULL REFERENCES group_cohorts(id) ON DELETE CASCADE;
	`

	if _, err := database.Exec(groupCohortIdAdd); err != nil {
		return fmt.Errorf("failed to add group_cohort_id column to student_groups: %w", err)
	}

	return nil
}
func groupMembersMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS group_members (
		id UUID PRIMARY KEY,
		student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
		student_group_id UUID REFERENCES student_groups(id) ON DELETE SET NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create group_members table: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "group_members"); err != nil {
		return fmt.Errorf("failed to create on update trigger for group_members: %w", err)
	}

	return nil
}
func teacherLoadMigrations(database *sqlx.DB) error {
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

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create teacher_loads table: %w", err)
	}

	createUniqueIndex := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_teacher_loads_unique
	ON teacher_loads (teacher_id, discipline_id, lesson_type_id);
	`
	if _, err := database.Exec(createUniqueIndex); err != nil {
		return fmt.Errorf("failed to create teacher_loads unique index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "teacher_loads"); err != nil {
		return fmt.Errorf("failed to create on update trigger for teacher_loads: %w", err)
	}

	// WARNING: it should be removed if group_cohort_id will be needed
	dropColumn := `
	ALTER TABLE teacher_loads
	DROP COLUMN IF EXISTS group_cohort_id;
	`

	if _, err := database.Exec(dropColumn); err != nil {
		return fmt.Errorf("failed to drop group_cohort_id column: %w", err)
	}

	return nil
}
func groupCohortMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS group_cohorts (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create group_cohorts table: %w", err)
	}

	indexSlug := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_group_cohorts_slug
	ON group_cohorts (slug);
	`
	if _, err := database.Exec(indexSlug); err != nil {
		return fmt.Errorf("failed to create group_cohorts slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "group_cohorts"); err != nil {
		return fmt.Errorf("failed to create on update trigger for group_cohorts: %w", err)
	}

	return nil
}
func groupCohortAssignmentMigrations(database *sqlx.DB) error {
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

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create group_cohort_assignments table: %w", err)
	}

	indexUnique := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_group_cohort_assignments_unique
	ON group_cohort_assignments (group_cohort_id, discipline_id, lesson_type_id);
	`
	if _, err := database.Exec(indexUnique); err != nil {
		return fmt.Errorf("failed to create unique index for group_cohort_assignments: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "group_cohort_assignments"); err != nil {
		return fmt.Errorf("failed to create on update trigger for group_cohort_assignments: %w", err)
	}

	return nil
}
