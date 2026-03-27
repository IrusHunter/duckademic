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
