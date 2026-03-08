package main

import (
	"fmt"
	"slices"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

// ==========================================================================================================
// =========================================== UpstreamRepository ===========================================
// ==========================================================================================================

type UpstreamRepository interface {
	Find(uuid.UUID) *Upstream
	FindFirstByName(string) *Upstream
	Add(Upstream) error
	Clear()
}

func NewUpstreamRepository(d *db.Database) UpstreamRepository {
	return &upstreamRepository{db: d}
}

type upstreamRepository struct {
	upstreams []Upstream
	db        *db.Database
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
func (r *upstreamRepository) Add(upstream Upstream) error {
	other := r.Find(upstream.ID)
	if other != nil {
		return fmt.Errorf("upstream %s already exists", upstream.Name)
	}

	r.upstreams = append(r.upstreams, upstream)
	return nil
}
func (r *upstreamRepository) Clear() {
	r.upstreams = []Upstream{}
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
