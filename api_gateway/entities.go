package main

import (
	"fmt"

	"github.com/google/uuid"
)

type Upstream struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	URL     string    `json:"url"`
	Enabled bool      `json:"enabled"`
}

func (u *Upstream) String() string {
	return fmt.Sprintf("id: %s, name: %s, url: %s, enabled: %v", u.ID, u.Name, u.URL, u.Enabled)
}

type Endpoint struct {
	Path         string   `json:"path"`
	UpstreamName string   `json:"upstream_name"`
	Upstream     Upstream `json:"-"`
}
