package main

import (
	"context"
	"fmt"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/jmoiron/sqlx"
)

func Migrate(database *sqlx.DB) error {
	migrationsF := []func(*sqlx.Tx) error{
		teacherMigrations,
		studentMigrations,
	}

	for _, f := range migrationsF {
		tx, err := database.Beginx()
		if err != nil {
			return err
		}

		if err := f(tx); err != nil {
			_ = tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func studentMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS students (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create students table: %w", err)
	}

	indexSlug := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_students_slug
	ON students (slug);
	`
	if _, err := tx.Exec(indexSlug); err != nil {
		return fmt.Errorf("failed to create students slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "students"); err != nil {
		return fmt.Errorf("failed to create on update trigger for students: %w", err)
	}

	return nil
}
func teacherMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS teachers (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create teachers table: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "teachers"); err != nil {
		return fmt.Errorf("failed to create on update trigger for teachers: %w", err)
	}

	return nil
}
