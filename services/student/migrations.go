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
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		middle_name TEXT,
		phone_number TEXT,
		email TEXT NOT NULL,
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
