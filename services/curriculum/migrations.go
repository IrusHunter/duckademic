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
		curriculumMigrations,
		semesterMigrations,
		lessonTypeMigrations,
		disciplineMigrations,
	}

	for _, f := range migrationsF {
		err := f(database)
		if err != nil {
			return err
		}
	}

	return nil
}

func curriculumMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS curriculums (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		duration_years INT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		effective_from TIMESTAMP WITH TIME ZONE NOT NULL,
		effective_to TIMESTAMP WITH TIME ZONE
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create curriculums table: %w", err)
	}

	indexSlug := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_curriculums_slug
	ON curriculums (slug);
	`
	if _, err := database.Exec(indexSlug); err != nil {
		return fmt.Errorf("failed to create curriculums slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "curriculums"); err != nil {
		return fmt.Errorf("failed to create on update trigger for curriculums: %w", err)
	}

	return nil
}
func semesterMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS semesters (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		curriculum_id UUID NOT NULL REFERENCES curriculums(id) ON DELETE CASCADE,
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
func lessonTypeMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS lesson_types (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		hours_value INT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create lesson_types table: %w", err)
	}

	indexSlug := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_lesson_types_slug
	ON lesson_types (slug);
	`
	if _, err := database.Exec(indexSlug); err != nil {
		return fmt.Errorf("failed to create lesson_types slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "lesson_types"); err != nil {
		return fmt.Errorf("failed to create on update trigger for lesson_types: %w", err)
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

	indexSlug := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_disciplines_slug
	ON disciplines (slug);
	`
	if _, err := database.Exec(indexSlug); err != nil {
		return fmt.Errorf("failed to create disciplines slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "disciplines"); err != nil {
		return fmt.Errorf("failed to create on update trigger for disciplines: %w", err)
	}

	return nil
}
