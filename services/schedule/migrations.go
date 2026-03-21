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

	return nil
}
