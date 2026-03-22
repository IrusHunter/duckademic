package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/student/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

func NewRESTAPI(
	sh resthandlers.StudentHandler,
	dh resthandlers.DatabaseHandler,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:   platform.NewRESTAPIHelper("RESTAPI"),
		studentHandler:  sh,
		databaseHandler: dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	studentHandler  resthandlers.StudentHandler
	databaseHandler resthandlers.DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/students", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.studentHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.studentHandler.Add),
	})
	ra.NewRoute("/student/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.studentHandler.Find),
		http.MethodDelete: ra.NewDefaultHandler(ra.studentHandler.Delete),
		http.MethodPut:    ra.NewDefaultHandler(ra.studentHandler.Update),
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
