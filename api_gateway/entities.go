package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ==========================================================================================================
// ================================================ Upstream ================================================
// ==========================================================================================================

// Upstream describes a backend service that requests can be proxied to.
type Upstream struct {
	ID        uuid.UUID `json:"id"`         // Unique identifier.
	Name      string    `json:"name"`       // Human-readable identifier.
	URL       string    `json:"url"`        // Base URL or domain.
	Enabled   bool      `json:"enabled"`    // Whether the service is enabled.
	CreatedAt time.Time `json:"created_at"` // Creation time.
	UpdatedAt time.Time `json:"updated_at"` // Last update time.
}

// String returns a human-readable representation of the Upstream.
// Includes ID, name, URL, enabled status, and optional created and updated timestamps.
func (u *Upstream) String() string {
	const timeFormat = "2006-01-02 15:04:05"

	var updatedAtStr, createdAtStr string
	if !u.CreatedAt.IsZero() {
		createdAtStr = fmt.Sprintf(", created: %s", u.CreatedAt.Format(timeFormat))
		updatedAtStr = fmt.Sprintf(", updated: %s", u.UpdatedAt.Format(timeFormat))
	}
	return fmt.Sprintf("Upstream{id: %s, name: %s, url: %s, enabled: %v%s%s}",
		u.ID, u.Name, u.URL, u.Enabled, createdAtStr, updatedAtStr,
	)
}

// ==========================================================================================================
// ================================================ Endpoint ================================================
// ==========================================================================================================

type Endpoint struct {
	Path         string   `json:"path"`
	UpstreamName string   `json:"upstream_name"`
	Upstream     Upstream `json:"-"`
}
