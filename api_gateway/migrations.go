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
	schema := ` 
	CREATE TABLE IF NOT EXISTS upstreams (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		enabled BOOLEAN NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create upstreams relation: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "upstreams"); err != nil {
		return fmt.Errorf("failed to create on update trigger for upstreams: %w", err)
	}

	schema = `
	CREATE TABLE IF NOT EXISTS endpoints (
		id UUID PRIMARY KEY,
		path TEXT NOT NULL UNIQUE,
		upstream_id UUID NOT NULL REFERENCES upstreams(id) ON DELETE CASCADE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create endpoints relation: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "endpoints"); err != nil {
		return fmt.Errorf("failed to create on update trigger for endpoints: %w", err)
	}

	return nil
}
