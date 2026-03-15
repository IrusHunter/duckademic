package platform

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RepositoryConfig struct {
	ClassName           string
	TableName           string
	EntityName          string
	AddParameters       []string
	GetParameters       []string
	UpdateParameters    []string
	ReturningParameters []string
}

// NewRepositoryConfig creates a new RepositoryConfig instance.
//
// It requires the class name (cn) table name (tn), entity name (en),
// and the lists of parameters for add (ap), get (gp), update (up), and returning (rp) operations.
func NewRepositoryConfig(
	cn string,
	tn string,
	en string,
	ap []string,
	gp []string,
	up []string,
	rp []string,
) RepositoryConfig {
	return RepositoryConfig{
		ClassName:           cn,
		TableName:           tn,
		EntityName:          en,
		AddParameters:       ap,
		GetParameters:       gp,
		UpdateParameters:    up,
		ReturningParameters: rp,
	}
}

// BaseRepository represents a base version for storing different entities.
type BaseRepository[T fmt.Stringer] interface {
	// Add inserts a new entity into the database and returns it, or an error if it fails.
	Add(context.Context, T) (T, error)
	Clear(context.Context) error // Clear removes all entities from the database.
	// FindByID returns a pointer to the entity with the given id from database.
	FindByID(context.Context, uuid.UUID) *T
	// FindFirstBy returns the first entity where the specified field matches the given value.
	FindFirstBy(ctx context.Context, field string, slug any) *T
	// GetAll returns a slice with all entities from database.
	GetAll(context.Context) []T
	// Delete removes the entity with the specified ID from the database.
	Delete(context.Context, uuid.UUID) error
	// Update updates the entity with the specified ID and returns the updated onr.
	Update(context.Context, uuid.UUID, T) (T, error)
	// SoftDelete marks the entity as deleted by setting the deleted_at timestamp.
	SoftDelete(context.Context, uuid.UUID) (T, error)
}

// NewBaseRepository creates a new BaseRepository instance.
//
// It requires a database connection (db) and a config (rc).
func NewBaseRepository[T fmt.Stringer](rc RepositoryConfig, db *sqlx.DB) BaseRepository[T] {
	return &baseRepository[T]{
		db:               db,
		logger:           logger.NewLogger(rc.ClassName+".txt", rc.ClassName),
		RepositoryConfig: rc,
	}
}

type baseRepository[T fmt.Stringer] struct {
	RepositoryConfig
	db        *sqlx.DB
	logger    logger.Logger
	nilEntity T
}

func (r *baseRepository[T]) Add(ctx context.Context, entity T) (T, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(%s)
		VALUES
		(%s)
		RETURNING %s
	`, r.TableName, r.FormSqlParameters(r.AddParameters),
		r.FormSqlValues(r.AddParameters), r.FormSqlParameters(r.ReturningParameters),
	)

	rows, err := r.db.NamedQueryContext(
		ctx,
		query,
		entity,
	)

	if err != nil {
		return r.nilEntity, r.logger.LogAndReturnError(
			contextutil.GetTraceID(ctx), "Add",
			fmt.Errorf("failed to insert %s: %w", entity.String(), err), logger.RepositoryQueryFailed,
		)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&entity); err != nil {
			return r.nilEntity, r.logger.LogAndReturnError(
				contextutil.GetTraceID(ctx), "Add",
				fmt.Errorf("failed to scan database row for %s: %w", entity.String(), err), logger.RepositoryScanFailed,
			)
		}
	} else {
		return r.nilEntity, r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Add",
			fmt.Errorf("no row returned after insert for %s", entity.String()), logger.RepositoryQueryFailed,
		)
	}

	r.logger.Log(contextutil.GetTraceID(ctx), "Add",
		fmt.Sprintf("%s successfully added", entity.String()), logger.RepositoryOperationSuccess)
	return entity, nil
}
func (r *baseRepository[T]) Clear(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, fmt.Sprintf(`DELETE FROM %s`, r.TableName))
	if err != nil {
		return r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Clear",
			fmt.Errorf("failed to truncate table %s: %w", r.TableName, err), logger.RepositoryQueryFailed,
		)
	}

	r.logger.Log(contextutil.GetTraceID(ctx), "Clear",
		fmt.Sprintf("table %s cleared", r.TableName),
		logger.RepositoryOperationSuccess)
	return nil
}
func (r *baseRepository[T]) FindByID(ctx context.Context, id uuid.UUID) *T {
	return r.FindFirstBy(ctx, "id", id)
}
func (r *baseRepository[T]) FindFirstBy(ctx context.Context, field string, param any) *T {
	query := fmt.Sprintf(
		`SELECT %s FROM %s
		WHERE %s=$1 LIMIT 1`,
		r.FormSqlParameters(r.GetParameters), r.TableName, field,
	)

	var entity T
	if err := r.db.GetContext(ctx, &entity, query, param); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			r.logger.Log(contextutil.GetTraceID(ctx), "FindFirstBy",
				fmt.Sprintf("%s with %s %q not found", r.EntityName, field, param),
				logger.RepositoryOperationSuccess,
			)
			return nil
		}
		r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "FindFirstBy",
			fmt.Errorf("failed to scan database row for %s %q: %w", field, param, err), logger.RepositoryScanFailed,
		)
		return nil
	}

	r.logger.Log(contextutil.GetTraceID(ctx), "FindFirstBy",
		fmt.Sprintf("%s found", entity.String()),
		logger.RepositoryOperationSuccess)
	return &entity
}
func (r *baseRepository[T]) GetAll(ctx context.Context) []T {
	query := fmt.Sprintf(`SELECT %s FROM %s`, r.FormSqlParameters(r.GetParameters), r.TableName)

	entities := []T{}
	err := r.db.SelectContext(ctx, &entities, query)
	if err != nil {
		r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "GetAll",
			fmt.Errorf("failed to get %ss: %w", r.EntityName, err), logger.RepositoryQueryFailed,
		)
	}

	r.logger.Log(contextutil.GetTraceID(ctx), "GetAll",
		fmt.Sprintf("%d entities found", len(entities)),
		logger.RepositoryOperationSuccess)
	return entities
}
func (r *baseRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, r.TableName)

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Delete",
			fmt.Errorf("failed to delete academic degree %q: %w", id, err), logger.RepositoryQueryFailed,
		)
	}

	r.logger.Log(contextutil.GetTraceID(ctx), "Delete",
		fmt.Sprintf("entity with id %q deleted", id),
		logger.RepositoryOperationSuccess)
	return nil
}
func (r *baseRepository[T]) Update(ctx context.Context, id uuid.UUID, entity T) (T, error) {
	query := fmt.Sprintf(`
		UPDATE %s SET
		%s
		WHERE id= :id
		RETURNING %s
		`, r.TableName, r.FormSqlEquations(r.UpdateParameters), r.FormSqlParameters(r.GetParameters))

	params := structToMapByDBTag(entity)
	params["id"] = id

	rows, err := r.db.NamedQueryContext(
		ctx,
		query,
		params,
	)
	if err != nil {
		return r.nilEntity, r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Update",
			fmt.Errorf("failed to update %s with id %q: %w", entity.String(), id, err), logger.RepositoryQueryFailed,
		)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&entity); err != nil {
			return r.nilEntity, r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Update",
				fmt.Errorf("failed to scan database row for %s with id %q: %w", entity.String(), id, err),
				logger.RepositoryScanFailed,
			)
		}
	} else {
		return r.nilEntity, r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Update",
			fmt.Errorf("%s with id %q not found to update", entity.String(), id), logger.RepositoryQueryFailed,
		)
	}

	r.logger.Log(contextutil.GetTraceID(ctx), "Update",
		fmt.Sprintf("%s successfully updated", entity.String()),
		logger.RepositoryOperationSuccess)
	return entity, nil
}
func (r *baseRepository[T]) SoftDelete(ctx context.Context, id uuid.UUID) (T, error) {
	if slices.Contains(r.GetParameters, "deleted_at") {
		return r.nilEntity, r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SoftDelete",
			fmt.Errorf("table %s does not support soft delete (missing deleted_at column)", r.TableName),
			logger.RepositoryQueryFailed,
		)
	}

	query := fmt.Sprintf(`
		UPDATE %s SET
			deleted_at = NOW()
		WHERE id = :id
		RETURNING %s
	`, r.TableName, r.FormSqlParameters(r.GetParameters))

	params := map[string]any{
		"id": id,
	}

	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return r.nilEntity, r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SoftDelete",
			fmt.Errorf("failed to soft delete %s with id %q: %w", r.TableName, id, err),
			logger.RepositoryQueryFailed,
		)
	}
	defer rows.Close()

	var entity T
	if rows.Next() {
		if err := rows.StructScan(&entity); err != nil {
			return r.nilEntity, r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SoftDelete",
				fmt.Errorf("failed to scan database row for %s with id %q: %w", r.TableName, id, err),
				logger.RepositoryScanFailed,
			)
		}
	} else {
		return r.nilEntity, r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "SoftDelete",
			fmt.Errorf("%s with id %q not found to delete", r.TableName, id),
			logger.RepositoryQueryFailed,
		)
	}

	r.logger.Log(contextutil.GetTraceID(ctx), "SoftDelete",
		fmt.Sprintf("%s with id %q successfully soft deleted", r.TableName, id),
		logger.RepositoryOperationSuccess)
	return entity, nil
}

func (r *baseRepository[T]) FormSqlParameters(parameters []string) string {
	return strings.Join(parameters, ", ")
}
func (r *baseRepository[T]) FormSqlValues(parameters []string) string {
	return ":" + strings.Join(parameters, ", :")
}
func (r *baseRepository[T]) FormSqlEquations(parameters []string) string {
	parts := make([]string, len(parameters))

	for i, p := range parameters {
		parts[i] = p + " = :" + p
	}

	return strings.Join(parts, ", ")
}

func structToMapByDBTag(entity interface{}) map[string]interface{} {
	v := reflect.ValueOf(entity)
	t := reflect.TypeOf(entity)

	// if v.Kind() == reflect.Ptr {
	// 	v = v.Elem()
	// 	t = t.Elem()
	// }

	result := make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			continue
		}
		result[dbTag] = v.Field(i).Interface()
	}

	return result
}
