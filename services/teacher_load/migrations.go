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
		teacherMigrations,
		disciplineMigrations,
		lessonTypeMigrations,
		teacherLoadMigrations,
		groupCohortMigrations,
	}

	for _, f := range migrationsF {
		err := f(database)
		if err != nil {
			return err
		}
	}

	return nil
}

func teacherMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS teachers (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		deleted_at TIMESTAMP WITH TIME ZONE
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create teachers table: %w", err)
	}

	indexName := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_teachers_slug
	ON teachers (slug);
	`

	if _, err := database.Exec(indexName); err != nil {
		return fmt.Errorf("failed to create teachers slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "teachers"); err != nil {
		return fmt.Errorf("failed to create on update trigger for teachers: %w", err)
	}

	dropDeletedAt := `
	ALTER TABLE teachers
	DROP COLUMN IF EXISTS deleted_at;
	`

	if _, err := database.Exec(dropDeletedAt); err != nil {
		return fmt.Errorf("failed to drop deleted_at column from teachers: %w", err)
	}

	return nil
}
func groupCohortMigrations(database *sqlx.DB) error {
	dropTable := `
	DROP TABLE IF EXISTS group_cohorts;
	`

	if _, err := database.Exec(dropTable); err != nil {
		return fmt.Errorf("failed to drop group_cohorts table: %w", err)
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
func teacherLoadMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS teacher_loads (
		id UUID PRIMARY KEY,
		teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
		discipline_id UUID NOT NULL REFERENCES disciplines(id) ON DELETE CASCADE,
		lesson_type_id UUID NOT NULL REFERENCES lesson_types(id) ON DELETE CASCADE,
		group_cohort_id UUID NOT NULL REFERENCES group_cohorts(id) ON DELETE CASCADE,
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
	ON teacher_loads (teacher_id, discipline_id, lesson_type_id, group_cohort_id);
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
