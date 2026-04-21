package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/IrusHunter/duckademic/shared/jsonutil"
)

// ==========================================================================================================
// ============================================== ProxyHandler ==============================================
// ==========================================================================================================

// DatabaseHandler represents a handler responsible for forwarding them to the appropriate upstream service.
type ProxyHandler interface {
	// HandlePath forwards request to the target upstream service.
	HandlePath(w http.ResponseWriter, r *http.Request)
}

// NewProxyHandler creates a new ProxyHandler instance.
//
// It requires a endpoint services (es) and an HTTP client (c).
func NewProxyHandler(es EndpointService, c *http.Client) ProxyHandler {
	return &proxyHandler{endpointService: es, client: c}
}

// proxyHandler is the basic implementation of the ProxyHandler interface.
type proxyHandler struct {
	endpointService EndpointService
	client          *http.Client
}

func (h *proxyHandler) HandlePath(w http.ResponseWriter, r *http.Request) {
	upstream, err := h.endpointService.GetRequiredService(r.URL.Path)
	if err != nil {
		jsonutil.ResponseWithError(w, http.StatusBadRequest, fmt.Errorf("server not found: %s", err))
		return
	}

	url := upstream.URL + "/" + strings.Join(strings.Split(r.URL.Path, "/")[2:], "/")
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		jsonutil.ResponseWithError(w, http.StatusInternalServerError,
			fmt.Errorf("failed to create request %q: %s", url, err.Error()),
		)
		return
	}

	for k, vv := range r.Header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}

	resp, err := h.client.Do(req)
	if err != nil {
		jsonutil.ResponseWithError(w, http.StatusInternalServerError, fmt.Errorf("request failed: %s", err.Error()))
		return
	}
	defer resp.Body.Close()

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Can't copy data")
	}
}

// ==========================================================================================================
// ============================================= DatabaseHandler ============================================
// ==========================================================================================================

// DatabaseHandler represents a handler responsible for database-related HTTP operations.
type DatabaseHandler interface {
	// Performs database seeding operations, initializing required data and the databases of the services via HTTP.
	Seed(w http.ResponseWriter, r *http.Request)
}

// NewDatabaseHandler creates a new DatabaseHandler instance.
//
// It requires the upstream (us) and endpoint (es) services, an HTTP client (c).
func NewDatabaseHandler(us UpstreamService, es EndpointService, c *http.Client) DatabaseHandler {
	return &databaseHandler{
		upstreamService: us,
		endpointService: es,
		client:          c,
	}
}

// databaseHandler is the basic implementation of the DatabaseHandler interface.
type databaseHandler struct {
	upstreamService UpstreamService
	endpointService EndpointService
	client          *http.Client
}

func (h *databaseHandler) Seed(w http.ResponseWriter, r *http.Request) {
	h.propagate(r, "/clear")

	if err := h.upstreamService.Seed(); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed upstreams: %w", err))
		return
	}
	if err := h.endpointService.Seed(); err != nil {
		jsonutil.ResponseWithError(w, 500, fmt.Errorf("failed to seed endpoints: %w", err))
		return
	}

	h.propagate(r, "/seed")

	jsonutil.ResponseWithJSON(w, 204, nil)
}

func (h *databaseHandler) propagate(r *http.Request, path string) {
	upstreams := h.upstreamService.GetAll(context.Background())

	for _, upstream := range upstreams {
		url := upstream.URL + path

		req, err := http.NewRequest(r.Method, url, r.Body)
		if err != nil {
			log.Printf("failed to create request %q: %s", url, err.Error())
			continue
		}

		resp, err := h.client.Do(req)
		if err != nil {
			log.Printf("request to %s failed: %s", url, err.Error())
			continue
		}

		respBody := map[string]string{}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		resp.Body.Close()
		if err != nil {
			continue
		}

		if respErr, ok := respBody["error"]; ok {
			log.Printf("can't call %s on %s: %s", path, upstream.String(), respErr)
		}
	}
}
