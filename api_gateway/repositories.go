package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// ==========================================================================================================
// =========================================== UpstreamRepository ===========================================
// ==========================================================================================================

// UpstreamRepository represents a storage for upstream service entities.
type UpstreamRepository interface {
	Find(context.Context, uuid.UUID) *Upstream // Find returns a pointer to the upstream with the given ID.
	// FindFirstByName returns a pointer to the first upstream with the given name.
	FindFirstByName(context.Context, string) *Upstream
	// Add inserts a new Upstream into the repository and returns it, or an error if it fails.
	Add(context.Context, Upstream) (Upstream, error)
	Clear(context.Context) // Clear removes all upstreams from the repository.
	// GetAll returns a slice with all upstreams from database.
	GetAll(context.Context) []Upstream
}

// NewUpstreamRepository creates a new UpstreamRepository instance.
//
// It requires a database connection (db).
func NewUpstreamRepository(d *sqlx.DB) UpstreamRepository {
	repo := &upstreamRepository{db: d}
	return repo
}

// upstreamRepository is the basic implementation of the UpstreamRepository interface.
type upstreamRepository struct {
	db *sqlx.DB
}

func (r *upstreamRepository) Find(ctx context.Context, id uuid.UUID) *Upstream {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, url, enabled, created_at, updated_at FROM upstreams
		WHERE id=$1`,
		id,
	)
	var upstream Upstream
	if err := row.Scan(
		&upstream.ID,
		&upstream.Name,
		&upstream.URL,
		&upstream.Enabled,
		&upstream.CreatedAt,
		&upstream.UpdatedAt,
	); err != nil {
		log.Printf("Can't scan upstreams row for id %q: %s \n", id, err.Error())
		return nil
	}

	return &upstream
}
func (r *upstreamRepository) FindFirstByName(ctx context.Context, name string) *Upstream {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, url, enabled, created_at, updated_at FROM upstreams
		WHERE name=$1 
		LIMIT 1`,
		name,
	)
	var upstream Upstream
	if err := row.Scan(
		&upstream.ID,
		&upstream.Name,
		&upstream.URL,
		&upstream.Enabled,
		&upstream.CreatedAt,
		&upstream.UpdatedAt,
	); err != nil {
		log.Printf("Can't scan upstreams row for name %q: %s \n", name, err.Error())
		return nil
	}

	return &upstream
}
func (r *upstreamRepository) Add(ctx context.Context, upstream Upstream) (Upstream, error) {
	rows, err := r.db.NamedQueryContext(
		ctx,
		`INSERT INTO upstreams 
		(id, name, url, enabled)
		VALUES
		(:id, :name, :url, :enabled)
		RETURNING created_at, updated_at`,
		upstream,
	)
	if err != nil {
		return Upstream{}, fmt.Errorf("failed to insert %q: %w", upstream.String(), err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&upstream.CreatedAt, &upstream.UpdatedAt); err != nil {
			return Upstream{}, fmt.Errorf("failed to scan database row for %s: %w", upstream.String(), err)
		}
	}

	return upstream, nil
}
func (r *upstreamRepository) Clear(ctx context.Context) {
	_, err := r.db.ExecContext(ctx, `DELETE FROM upstreams`)
	if err != nil {
		log.Println("Can't truncate table upstreams: " + err.Error())
	}
}
func (r *upstreamRepository) GetAll(ctx context.Context) []Upstream {
	upstreams := []Upstream{}
	err := r.db.SelectContext(
		ctx,
		&upstreams,
		`SELECT id, name, url, enabled, created_at, updated_at FROM upstreams`,
	)

	if err != nil {
		log.Println(fmt.Errorf("failed to get upstreams: %w", err).Error())
	}

	return upstreams
}

// ==========================================================================================================
// =========================================== EndpointRepository ===========================================
// ==========================================================================================================

// EndpointRepository represents a storage for endpoint entities.
type EndpointRepository interface {
	// Add inserts a new Endpoint into the repository and returns it, or an error if it fails.
	Add(context.Context, Endpoint) (Endpoint, error)
	Clear(context.Context) // Clear removes all endpoints from the repository.
	// Refresh reloads endpoints from the underlying storage, returning an error on failure.
	Refresh(context.Context) error
	// FindFirstByName returns a pointer to the endpoint with the given path.
	FindByPath(string) *Endpoint
}

// NewEndpointRepository creates a new EndpointRepository instance.
//
// It requires a database connection (db) and upstream repository (ur).
func NewEndpointRepository(db *sqlx.DB, ur UpstreamRepository) EndpointRepository {
	repo := &endpointRepository{
		upstreamRepository: ur,
		endpoints:          map[string]Endpoint{},
		db:                 db,
	}
	if err := repo.Refresh(context.Background()); err != nil {
		log.Println("Can't refresh upstream repository: " + err.Error())
	}
	return repo
}

// endpointRepository is the basic implementation of the EndpointRepository interface.
type endpointRepository struct {
	endpoints          map[string]Endpoint
	upstreamRepository UpstreamRepository
	db                 *sqlx.DB
}

func (r *endpointRepository) FindByPath(path string) *Endpoint {
	endpoint, ok := r.endpoints[path]

	if !ok {
		return nil
	}
	return &endpoint
}
func (r *endpointRepository) Add(ctx context.Context, endpoint Endpoint) (Endpoint, error) {
	rows, err := r.db.NamedQueryContext(
		ctx,
		`INSERT INTO endpoints
		(id, path, upstream_id)
		VALUES
		(:id, :path, :upstream_id)
		Returning created_at, updated_at`,
		endpoint,
	)
	if err != nil {
		return Endpoint{}, fmt.Errorf("failed to insert endpoint %q: %w", endpoint.String(), err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&endpoint.CreatedAt, &endpoint.UpdatedAt); err != nil {
			return Endpoint{}, fmt.Errorf("failed to scan database row for endpoint %s: %w", endpoint.String(), err)
		}
	}

	if endpoint.Upstream == nil {
		upstream := r.upstreamRepository.Find(ctx, endpoint.UpstreamID)
		if upstream == nil {
			return Endpoint{}, fmt.Errorf("failed to find upstream for %s", endpoint.String())
		}
		endpoint.Upstream = upstream
	}

	r.endpoints[endpoint.Path] = endpoint
	return r.endpoints[endpoint.Path], nil
}
func (r *endpointRepository) Clear(ctx context.Context) {
	_, err := r.db.ExecContext(ctx, `DELETE FROM endpoints`)
	if err != nil {
		log.Println("Can't truncate table endpoints: " + err.Error())
	}
	r.endpoints = make(map[string]Endpoint)
}
func (r *endpointRepository) Refresh(ctx context.Context) error {
	type endpointWithUpstream struct {
		Endpoint
		UpstreamID      uuid.UUID `db:"upstream_id"`
		UpstreamName    string    `db:"upstream_name"`
		UpstreamURL     string    `db:"upstream_url"`
		UpstreamEnabled bool      `db:"upstream_enabled"`
	}
	var rows []endpointWithUpstream

	query := `
		SELECT e.id, e.path, e.upstream_id, e.created_at, e.updated_at, u.id AS upstream_id, 
			u.name AS upstream_name, u.url AS upstream_url, u.enabled AS upstream_enabled
		FROM endpoints e
		JOIN upstreams u ON u.id = e.upstream_id
	`
	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return fmt.Errorf("failed to get endpoints with upstreams: %w", err)
	}

	r.endpoints = make(map[string]Endpoint, len(rows))
	for _, row := range rows {
		upstream := &Upstream{
			ID:      row.UpstreamID,
			Name:    row.UpstreamName,
			URL:     row.UpstreamURL,
			Enabled: row.UpstreamEnabled,
		}

		endpoint := row.Endpoint
		endpoint.Upstream = upstream

		r.endpoints[endpoint.Path] = endpoint
	}

	return nil
}
