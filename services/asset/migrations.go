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
		classroomMigrations,
	}

	for _, f := range migrationsF {
		err := f(database)
		if err != nil {
			return err
		}
	}

	return nil
}

func classroomMigrations(database *sqlx.DB) error {
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

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create classrooms table: %w", err)
	}

	indexSlug := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_classrooms_slug
	ON classrooms (slug);
	`
	if _, err := database.Exec(indexSlug); err != nil {
		return fmt.Errorf("failed to create classrooms slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "classrooms"); err != nil {
		return fmt.Errorf("failed to create on update trigger for classrooms: %w", err)
	}

	return nil
}
