package main

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// ==========================================================================================================
// =========================================== UpstreamRepository ===========================================
// ==========================================================================================================

// UpstreamRepository represents a storage for upstream service entities.
type UpstreamRepository interface {
	Find(uuid.UUID) *Upstream         // Returns a pointer to the upstream with the given ID.
	FindFirstByName(string) *Upstream // Returns a pointer to the first upstream with the given name.
	// Add inserts a new Upstream into the repository and returns it, or an error if it fails.
	Add(context.Context, Upstream) (Upstream, error)
	Clear(context.Context) // Clear removes all upstreams from the repository.
	Refresh() error        // Refresh reloads upstreams from the underlying storage, returning an error on failure.
}

// NewUpstreamRepository creates a new UpstreamRepository instance.
//
// It requires a database connection (db).
func NewUpstreamRepository(d *sqlx.DB) UpstreamRepository {
	return &upstreamRepository{db: d}
}

// upstreamRepository is the basic implementation of the UpstreamRepository interface.
type upstreamRepository struct {
	upstreams []Upstream
	db        *sqlx.DB
}

func (r *upstreamRepository) Find(id uuid.UUID) *Upstream {
	ind := slices.IndexFunc(r.upstreams, func(other Upstream) bool {
		return other.ID == id
	})

	if ind == -1 {
		return nil
	}
	return &r.upstreams[ind]
}
func (r *upstreamRepository) FindFirstByName(name string) *Upstream {
	ind := slices.IndexFunc(r.upstreams, func(other Upstream) bool {
		return other.Name == name
	})

	if ind == -1 {
		return nil
	}
	return &r.upstreams[ind]
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
		return Upstream{}, fmt.Errorf("failed to insert upstream %q: %w", upstream.String(), err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&upstream.CreatedAt, &upstream.UpdatedAt); err != nil {
			return Upstream{}, err
		}
	}

	r.upstreams = append(r.upstreams, upstream)
	return upstream, nil
}
func (r *upstreamRepository) Clear(ctx context.Context) {
	r.db.ExecContext(ctx, `TRUNCATE TABLE upstreams`)
	r.upstreams = []Upstream{}
}
func (r *upstreamRepository) Refresh() error {
	return nil
}

// ==========================================================================================================
// =========================================== EndpointRepository ===========================================
// ==========================================================================================================

type EndpointRepository interface {
	Add(Endpoint) error
	Find(string) *Endpoint
}

func NewEndpointRepository(ur UpstreamRepository) EndpointRepository {
	return &endpointRepository{upstreamRepository: ur, endpoints: map[string]Endpoint{}}
}

type endpointRepository struct {
	endpoints          map[string]Endpoint
	upstreamRepository UpstreamRepository
}

func (r *endpointRepository) Find(path string) *Endpoint {
	endpoint, ok := r.endpoints[path]

	if !ok {
		return nil
	}
	return &endpoint
}
func (r *endpointRepository) Add(endpoint Endpoint) error {
	other := r.Find(endpoint.Path)
	if other != nil {
		return fmt.Errorf("endpoint with path %s already exists", endpoint.Path)
	}

	upstream := r.upstreamRepository.FindFirstByName(endpoint.UpstreamName)
	if upstream == nil {
		return fmt.Errorf("upstream %s not found", endpoint.UpstreamName)
	}
	endpoint.Upstream = *upstream

	r.endpoints[endpoint.Path] = endpoint
	return nil
}
