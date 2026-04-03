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
		semesterMigrations,
		studentMigrations,
		groupCohortMigrations,
		studentGroupMigrations,
		groupMembersMigrations,
		disciplineMigrations,
		lessonTypeMigrations,
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

func semesterMigrations(database *sqlx.DB) error {
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

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create semesters table: %w", err)
	}

	indexSlug := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_semesters_slug
	ON semesters (slug);
	`
	if _, err := database.Exec(indexSlug); err != nil {
		return fmt.Errorf("failed to create semesters slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "semesters"); err != nil {
		return fmt.Errorf("failed to create on update trigger for semesters: %w", err)
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

	addColumnWithFK := `
	ALTER TABLE students
	ADD COLUMN IF NOT EXISTS semester_id UUID REFERENCES semesters(id) ON DELETE SET NULL;
	`

	if _, err := database.Exec(addColumnWithFK); err != nil {
		return fmt.Errorf("failed to add semester_id column with foreign key to students: %w", err)
	}

	return nil
}
func groupCohortMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS group_cohorts (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		semester_id UUID NOT NULL REFERENCES semesters(id) ON DELETE CASCADE,
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
func studentGroupMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS student_groups (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		group_cohort_id UUID NOT NULL REFERENCES group_cohorts(id) ON DELETE CASCADE,
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

	return nil
}
func groupMembersMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS group_members (
		id UUID PRIMARY KEY,
		student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
		group_cohort_id UUID NOT NULL REFERENCES group_cohorts(id) ON DELETE CASCADE,
		student_group_id UUID REFERENCES student_groups(id) ON DELETE SET NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		UNIQUE (student_id, group_cohort_id)
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create group_members table: %w", err)
	}

	indexStudentCohort := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_group_members_student_cohort
	ON group_members (student_id, group_cohort_id);
	`
	if _, err := database.Exec(indexStudentCohort); err != nil {
		return fmt.Errorf("failed to create group_members unique index on student_id and group_cohort_id: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "group_members"); err != nil {
		return fmt.Errorf("failed to create on update trigger for group_members: %w", err)
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
		return fmt.Errorf("failed to create lesson_types slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "lesson_types"); err != nil {
		return fmt.Errorf("failed to create on update trigger for lesson_types: %w", err)
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
