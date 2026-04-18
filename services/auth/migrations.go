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
		permissionMigrations,
	}

	for _, f := range migrationsF {
		err := f(database)
		if err != nil {
			return err
		}
	}

	return nil
}

func permissionMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS permissions (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create permissions table: %w", err)
	}

	indexName := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_permissions_name
	ON permissions (name);
	`

	if _, err := database.Exec(indexName); err != nil {
		return fmt.Errorf("failed to create permissions name index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "permissions"); err != nil {
		return fmt.Errorf("failed to create on update trigger for permissions: %w", err)
	}

	return nil
}
