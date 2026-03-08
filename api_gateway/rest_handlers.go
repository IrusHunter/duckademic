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

func (c *proxyHandler) HandlePath(w http.ResponseWriter, r *http.Request) {
	upstream, err := c.endpointService.GetRequiredService(r.URL.Path)
	if err != nil {
		jsonutil.ResponseWithError(w, http.StatusBadRequest, fmt.Errorf("server not found: %s", err))
		return
	}

	url := upstream.URL + r.URL.Path
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}
	req, err := http.NewRequest(r.Method, url, r.Body)

	resp, err := c.client.Do(req)
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
// =============================================== SeedHandler ==============================================
// ==========================================================================================================
