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

	return nil
}
