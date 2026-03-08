package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// EnsureUpdatedAtTrigger creates a trigger for the given table that
// automatically updates the `updated_at` column on row updates.
//
// It requires a context (ctx), a sqlx.DB connection (db), and the table name (tableName) to create the trigger for.
//
// Returns an error if the function or trigger cannot be created.
func EnsureUpdatedAtTrigger(ctx context.Context, db *sqlx.DB, tableName string) error {
	query := fmt.Sprintf(`
	CREATE OR REPLACE FUNCTION set_updated_at()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = NOW();
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	DO $$
	BEGIN
		IF NOT EXISTS (
			SELECT 1
			FROM pg_trigger
			WHERE tgname = 'trg_%s_updated_at'
		) THEN
			CREATE TRIGGER trg_%s_updated_at
			BEFORE UPDATE ON %s
			FOR EACH ROW
			EXECUTE FUNCTION set_updated_at();
		END IF;
	END
	$$;
	`, tableName, tableName, tableName,
	)

	_, err := db.ExecContext(ctx, query)
	return err
}
