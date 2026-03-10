package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/google/uuid"
)

// ==========================================================================================================
// ============================================= UpstreamService ============================================
// ==========================================================================================================

// UpstreamService provides operations to initialize and manage upstreams.
type UpstreamService interface {
	Seed() error // Seed clears existing upstream data and initializes it from a JSON file.
	// GetAll returns a slice with all upstreams from repository.
	GetAll(context.Context) []Upstream
}

// NewUpstreamService creates a new UpstreamService instance.
//
// It requires a upstream repository (ur).
func NewUpstreamService(ur UpstreamRepository) UpstreamService {
	return &upstreamService{repository: ur}
}

// upstreamService is the basic implementation of the UpstreamService interface.
type upstreamService struct {
	repository UpstreamRepository
}

func (s *upstreamService) Seed() error {
	upstreams := []Upstream{}
	if err := jsonutil.ReadFileTo(filepath.Join("data", "upstreams.json"), &upstreams); err != nil {
		return fmt.Errorf("failed to load upstreams seed data: %w", err)
	}

	s.repository.Clear(context.Background())
	for _, upstream := range upstreams {
		upstream.ID = uuid.New()
		_, err := s.repository.Add(context.Background(), upstream)
		if err != nil {
			s.repository.Clear(context.Background())
			return fmt.Errorf("can't add upstream %q: %w", upstream.String(), err)
		}
	}

	return nil
}
func (s *upstreamService) GetAll(ctx context.Context) []Upstream {
	return s.repository.GetAll(ctx)
}

// ==========================================================================================================
// ============================================= EndpointService ============================================
// ==========================================================================================================

// EndpointService provides operations to initialize and manage endpoints.
type EndpointService interface {
	// GetRequiredService returns the upstream service responsible for the given request path.
	GetRequiredService(string) (*Upstream, error)
	Seed() error // Seed clears existing endpoints data and initializes it from a JSON file.
}

// NewEndpointService creates a new EndpointService instance.
//
// It requires the endpoint (er) and upstream (ur) repositories.
func NewEndpointService(er EndpointRepository, ur UpstreamRepository) EndpointService {
	return &endpointService{repository: er, upstreamRepository: ur}
}

type endpointService struct {
	repository         EndpointRepository
	upstreamRepository UpstreamRepository
}

func (s *endpointService) GetRequiredService(path string) (*Upstream, error) {
	nPath, err := s.normalizePath(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path %q: %s", path, err.Error())
	}

	endpoint := s.repository.FindByPath(nPath)
	if endpoint == nil {
		return nil, fmt.Errorf("endpoint with path %q not found", nPath)
	}
	return endpoint.Upstream, nil
}
func (s *endpointService) normalizePath(path string) (string, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("only %d parts, should be more than 1", len(parts))
	}
	return strings.Split(parts[1], "?")[0], nil
}

func (s *endpointService) Seed() error {
	type jsonEndpoint struct {
		Path         string `json:"path"`
		UpstreamName string `json:"upstream_name"`
	}
	endpoints := []jsonEndpoint{}
	if err := jsonutil.ReadFileTo(filepath.Join("data", "endpoints.json"), &endpoints); err != nil {
		return fmt.Errorf("failed to load endpoints seed data: %w", err)
	}

	s.repository.Clear(context.Background())
	for _, je := range endpoints {
		upstream := s.upstreamRepository.FindFirstByName(context.Background(), je.UpstreamName)
		if upstream == nil {
			return fmt.Errorf("upstream with name %s not found", je.UpstreamName)
		}
		endpoint := Endpoint{
			ID:         uuid.New(),
			Path:       je.Path,
			UpstreamID: upstream.ID,
			Upstream:   upstream,
		}
		_, err := s.repository.Add(context.Background(), endpoint)
		if err != nil {
			s.repository.Clear(context.Background())
			return fmt.Errorf("failed to add %s to repository: %w", endpoint.String(), err)
		}
	}

	return nil
}
