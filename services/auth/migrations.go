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
		roleMigrations,
		rolePermissionsMigrations,
		serviceMigrations,
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
func roleMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS roles (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create roles table: %w", err)
	}

	indexName := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_name
	ON roles (name);
	`

	if _, err := database.Exec(indexName); err != nil {
		return fmt.Errorf("failed to create roles name index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "roles"); err != nil {
		return fmt.Errorf("failed to create on update trigger for roles: %w", err)
	}

	return nil
}
func rolePermissionsMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS role_permissions (
		id UUID PRIMARY KEY,
		role_id UUID NOT NULL,
		permission_id UUID NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create role_permissions table: %w", err)
	}

	if err := db.EnsureForeignKey(
		context.Background(),
		database,
		"fk_role_permissions_role",
		"role_permissions",
		"role_id",
		"roles",
		"id",
	); err != nil {
		return fmt.Errorf("failed to create role_id foreign key: %w", err)
	}

	if err := db.EnsureForeignKey(
		context.Background(),
		database,
		"fk_role_permissions_permission",
		"role_permissions",
		"permission_id",
		"permissions",
		"id",
	); err != nil {
		return fmt.Errorf("failed to create permission_id foreign key: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "role_permissions"); err != nil {
		return fmt.Errorf("failed to create on update trigger for role_permissions: %w", err)
	}

	return nil
}
func serviceMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS services (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		secrete TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create services table: %w", err)
	}

	indexName := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_services_name
	ON services (name);
	`

	if _, err := database.Exec(indexName); err != nil {
		return fmt.Errorf("failed to create services name index: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "services"); err != nil {
		return fmt.Errorf("failed to create on update trigger for services: %w", err)
	}

	return nil
}
func servicePermissionsMigrations(database *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS service_permissions (
		id UUID PRIMARY KEY,
		service_id UUID NOT NULL,
		permission_id UUID NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create service_permissions table: %w", err)
	}

	if err := db.EnsureForeignKey(
		context.Background(),
		database,
		"fk_service_permissions_service",
		"service_permissions",
		"service_id",
		"services",
		"id",
	); err != nil {
		return fmt.Errorf("failed to create service_id foreign key: %w", err)
	}

	if err := db.EnsureForeignKey(
		context.Background(),
		database,
		"fk_service_permissions_permission",
		"service_permissions",
		"permission_id",
		"permissions",
		"id",
	); err != nil {
		return fmt.Errorf("failed to create permission_id foreign key: %w", err)
	}

	if err := db.EnsureUpdatedAtTrigger(context.Background(), database, "service_permissions"); err != nil {
		return fmt.Errorf("failed to create on update trigger for service_permissions: %w", err)
	}

	return nil
}
