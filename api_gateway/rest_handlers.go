package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/IrusHunter/duckademic/shared/jsonutil"
)

// ==========================================================================================================
// ============================================== ProxyHandler ==============================================
// ==========================================================================================================

type ProxyHandler interface {
	HandlePath(w http.ResponseWriter, r *http.Request)
}

func NewProxyHandler(es EndpointService, c *http.Client) ProxyHandler {
	return &proxyHandler{endpointService: es, client: c}
}

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

	url := upstream.URL + r.URL.Path
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}
	req, err := http.NewRequest(r.Method, url, r.Body)

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
// It requires a upstream service (us).
func NewDatabaseHandler(us UpstreamService) DatabaseHandler {
	return &databaseHandler{
		upstreamService: us,
	}
}

// databaseHandler is the basic implementation of the DatabaseHandler interface.
type databaseHandler struct {
	upstreamService UpstreamService
}

func (h *databaseHandler) Seed(w http.ResponseWriter, r *http.Request) {
	err := h.upstreamService.Seed()
	if err != nil {
		jsonutil.ResponseWithError(w, 500, err)
		return
	}
	jsonutil.ResponseWithJSON(w, 204, nil)
}
