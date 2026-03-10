package main

import (
	"log"
	"net/http"
	"strconv"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

// NewRESTAPI creates a new RESTAPI instance.
//
// It requires the proxy (ph) and database (dh) handlers.
func NewRESTAPI(ph ProxyHandler, dh DatabaseHandler) RESTAPI {
	return &restapi{proxyHandler: ph, databaseHandler: dh}
}

// restapi is the basic implementation of the RESTAPI interface.
type restapi struct {
	proxyHandler    ProxyHandler
	databaseHandler DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	newHandler("/", ra.proxyHandler.HandlePath)
	newHandler("/seed", ra.databaseHandler.Seed)

	log.Printf("Server start at port %d \n", port)

	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func newHandler(path string, f http.HandlerFunc) {
	http.HandleFunc(path, corsMiddleware(f))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
