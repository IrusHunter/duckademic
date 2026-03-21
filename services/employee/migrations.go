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
		academicDegreeMigrations,
		employeeMigrations,
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
func academicDegreeMigrations(database *sqlx.DB) error {
	schema := ` 
	CREATE TABLE IF NOT EXISTS academic_degrees (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		title TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create academic degrees relation: %w", err)
	}

	index := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_academic_degrees_slug
	ON academic_degrees (slug);
	`

	if _, err := database.Exec(index); err != nil {
		return fmt.Errorf("failed to create academic_degrees slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "academic_degrees"); err != nil {
		return fmt.Errorf("failed to create on update trigger for academic_degrees: %w", err)
	}

	return nil
}
func employeeMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS employees (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		middle_name TEXT,
		phone_number TEXT,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		deleted_at TIMESTAMP WITH TIME ZONE
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create employees table: %w", err)
	}

	index := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_employees_slug
	ON employees (slug);
	`

	if _, err := database.Exec(index); err != nil {
		return fmt.Errorf("failed to create employees slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "employees"); err != nil {
		return fmt.Errorf("failed to create on update trigger for employees: %w", err)
	}

	return nil
}
func teacherMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS teachers (
		employee_id UUID PRIMARY KEY,
		email TEXT NOT NULL,
		academic_degree_id UUID,
		academic_rank_id UUID NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		deleted_at TIMESTAMP WITH TIME ZONE,

		CONSTRAINT fk_teachers_employee
			FOREIGN KEY (employee_id)
			REFERENCES employees(id)
			ON DELETE CASCADE,

		CONSTRAINT fk_teachers_academic_degree
			FOREIGN KEY (academic_degree_id)
			REFERENCES academic_degrees(id)
			ON DELETE SET NULL,

		CONSTRAINT fk_teachers_academic_rank
			FOREIGN KEY (academic_rank_id)
			REFERENCES academic_ranks(id)
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create teachers table: %w", err)
	}

	index := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_teachers_email
	ON teachers (email);
	`

	if _, err := database.Exec(index); err != nil {
		return fmt.Errorf("failed to create teachers email index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "teachers"); err != nil {
		return fmt.Errorf("failed to create on update trigger for teachers: %w", err)
	}

	return nil
}
