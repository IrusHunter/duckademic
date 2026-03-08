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

// UpstreamService provides operations to initialize and manage upstream services.
type UpstreamService interface {
	Seed() error // clears existing upstream data and initializes it from a JSON file.
}

// NewUpstreamService creates a new UpstreamService instance.
//
// It requires a upstream repository (up).
func NewUpstreamService(up UpstreamRepository) UpstreamService {
	return &upstreamService{repository: up}
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

// ==========================================================================================================
// ============================================= EndpointService ============================================
// ==========================================================================================================

type EndpointService interface {
	GetRequiredService(string) (Upstream, error)
}

func NewEndpointService(er EndpointRepository) EndpointService {
	return &endpointService{repository: er}
}

type endpointService struct {
	repository EndpointRepository
}

func (s *endpointService) GetRequiredService(path string) (Upstream, error) {
	nPath, err := s.normalizePath(path)
	if err != nil {
		return Upstream{}, fmt.Errorf("invalid path (%s): %s", path, err.Error())
	}

	endpoint := s.repository.Find(nPath)
	if endpoint == nil {
		return Upstream{}, fmt.Errorf("endpoint %s not found", nPath)
	}
	return endpoint.Upstream, nil
}
func (s *endpointService) normalizePath(path string) (string, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("only %d parts, should be more than 2", len(parts))
	}
	return strings.Split(parts[1], "?")[0], nil
}
