package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/auth/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

func NewRESTAPI(
	ph resthandlers.PermissionHandler,
	dh resthandlers.DatabaseHandler,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:     platform.NewRESTAPIHelper("RESTAPI"),
		permissionHandler: ph,
		databaseHandler:   dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	permissionHandler resthandlers.PermissionHandler
	databaseHandler   resthandlers.DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/permissions", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.permissionHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.permissionHandler.Add),
	})
	ra.NewRoute("/permission/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.permissionHandler.Find),
		http.MethodDelete: ra.NewDefaultHandler(ra.permissionHandler.Delete),
		http.MethodPut:    ra.NewDefaultHandler(ra.permissionHandler.Update),
	})

	http.HandleFunc("/seed", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.databaseHandler.Seed)(r.Context(), w, r)
	})
	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.databaseHandler.Clear)(r.Context(), w, r)
	})

	log.Printf("Server start at port %d \n", port)

	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
