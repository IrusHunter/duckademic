package db

import (
	"fmt"

	"github.com/IrusHunter/duckademic/shared/envutil"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Database wraps a sqlx.DB client and provides a centralized
// way to access the database throughout the application.
type Database struct {
	Client *sqlx.DB // The underlying database client used for queries.
}

// NewDatabase creates a new Database instance. If the specified database does not exist, it will be created.
//
// It requires host, port, user, password, dbname, sslmode for connecting to PostgreSQL.
//
// Returns a Database instance or an error if the connection or database creation fails.
func NewDatabase(host, port, user, dbname, password, sslmode string) (*Database, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, user, dbname, password, sslmode,
	)

	dbConn, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "3D000" {
			// log.Printf("Database %s does not exist, creating...", dbname)

			adminConnStr := fmt.Sprintf(
				"host=%s port=%s user=%s dbname=postgres password=%s sslmode=%s",
				host, port, user, password, sslmode,
			)
			adminDB, err := sqlx.Connect("postgres", adminConnStr)
			if err != nil {
				return nil, fmt.Errorf("failed to connect to admin database: %w", err)
			}
			defer adminDB.Close()

			_, err = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %q", dbname))
			if err != nil {
				return nil, fmt.Errorf("failed to create database %s: %w", dbname, err)
			}

			dbConn, err = sqlx.Connect("postgres", connectionString)
			if err != nil {
				return nil, fmt.Errorf("failed to connect to newly created database %s: %w", dbname, err)
			}
		} else {
			return nil, fmt.Errorf("failed to connect to %s: %w", dbname, err)
		}
	}
	return &Database{Client: dbConn}, nil
}

// NewDefaultConnection creates a new Database instance using the NewDatabase function.
// It reads the required configuration from environment variables.
//
// Required environment variables: DB_HOST, DB_PORT, DB_USER, DB_NAME, DB_PASSWORD, DB_SSLMODE.
//
// Returns a Database instance or an error if the connection fails.
func NewDefaultConnection() (*Database, error) {
	return NewDatabase(
		envutil.GetStringFromENV("DB_HOST"),
		envutil.GetStringFromENV("DB_PORT"),
		envutil.GetStringFromENV("DB_USER"),
		envutil.GetStringFromENV("DB_NAME"),
		envutil.GetStringFromENV("DB_PASSWORD"),
		envutil.GetStringFromENV("DB_SSLMODE"),
	)
}
