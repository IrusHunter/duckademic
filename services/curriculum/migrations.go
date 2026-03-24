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
