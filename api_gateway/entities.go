package main

import (
	"fmt"
	"time"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

// ==========================================================================================================
// ================================================ Upstream ================================================
// ==========================================================================================================

// Upstream represents a backend service that requests can be proxied to.
type Upstream struct {
	ID        uuid.UUID `db:"id" json:"id"`                 // Unique identifier.
	Name      string    `db:"name" json:"name"`             // Human-readable identifier.
	URL       string    `db:"url" json:"url"`               // Base URL or domain.
	Enabled   bool      `db:"enabled" json:"enabled"`       // Whether the service is enabled.
	CreatedAt time.Time `db:"created_at" json:"created_at"` // Record creation timestamp.
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"` // Record last update timestamp.
}

// String returns a human-readable representation of the Upstream.
// Includes ID, name, URL, enabled status, and optional created and updated timestamps.
func (u *Upstream) String() string {
	var updatedAtStr, createdAtStr string
	if !u.CreatedAt.IsZero() {
		createdAtStr = fmt.Sprintf(", created: %s", u.CreatedAt.Format(db.TimeFormat))
		updatedAtStr = fmt.Sprintf(", updated: %s", u.UpdatedAt.Format(db.TimeFormat))
	}
	return fmt.Sprintf("Upstream{id: %s, name: %s, url: %s, enabled: %v%s%s}",
		u.ID, u.Name, u.URL, u.Enabled, createdAtStr, updatedAtStr,
	)
}

// ==========================================================================================================
// ================================================ Endpoint ================================================
// ==========================================================================================================

// Endpoint represents a request path that should be proxied to an upstream service.
type Endpoint struct {
	ID         uuid.UUID `db:"id" json:"id"`                   // Unique identifier.
	Path       string    `db:"path" json:"path"`               // First url path segment used for routing (must be unique).
	UpstreamID uuid.UUID `db:"upstream_id" json:"upstream_id"` // ID of the upstream service that handles this endpoint.
	CreatedAt  time.Time `db:"created_at" json:"created_at"`   // Record creation timestamp.
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`   // Record last update timestamp.
	Upstream   *Upstream `db:"-" json:"-"`                     // Cached upstream instance used at runtime for fast routing.
}

// String returns a human-readable representation of the Endpoint.
// Includes ID, path, id of the upstream, and optional created and updated timestamps.
func (e *Endpoint) String() string {
	var createdAtStr, updatedAtStr string
	if !e.CreatedAt.IsZero() {
		createdAtStr = fmt.Sprintf("created_at: %s", e.CreatedAt.Format(db.TimeFormat))
		createdAtStr = fmt.Sprintf("created_at: %s", e.CreatedAt.Format(db.TimeFormat))
	}
	return fmt.Sprintf("Endpoint{id: %s, path: %s, UpstreamID: %s%s%s}",
		e.ID, e.Path, e.UpstreamID, createdAtStr, updatedAtStr,
	)
}
