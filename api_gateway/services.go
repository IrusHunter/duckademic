package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/IrusHunter/duckademic/shared/jsonutil"
)

// ==========================================================================================================
// ============================================= UpstreamService ============================================
// ==========================================================================================================

type UpstreamService interface {
	Seed() error
}

func NewUpstreamService(up UpstreamRepository) UpstreamService {
	return &upstreamService{repository: up}
}

type upstreamService struct {
	repository UpstreamRepository
}

func (s *upstreamService) Seed() error {
	upstreams := []Upstream{}
	if err := jsonutil.ReadFileTo(filepath.Join("data", "upstreams"), &upstreams); err != nil {
		return fmt.Errorf("failed to load upstreams seed data: %w", err)
	}

	for _, upstream := range upstreams {
		err := s.repository.Add(upstream)
		if err != nil {
			s.repository.Clear()
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
	dat, err := os.ReadFile(filepath.Join("data", "endpoints.json"))
	if err != nil {
		panic("can't read file data/endpoints.json: " + err.Error())
	}

	endpoints := make([]Endpoint, 0)
	err = json.Unmarshal(dat, &endpoints)
	if err != nil {
		panic("can't unmarshal data: " + err.Error())
	}

	for _, endpoint := range endpoints {
		err := er.Add(endpoint)
		if err != nil {
			log.Fatalf("can't add endpoint with path %s to collection: %s", endpoint.Path, err.Error())
		}
	}

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
