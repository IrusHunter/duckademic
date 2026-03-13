package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/employees/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

// NewRESTAPI creates a new RESTAPI instance.
//
// It requires the academic rank (arh)? academic degree (adh), and database (dh) handlers.
func NewRESTAPI(
	arh resthandlers.AcademicRankHandler,
	adh resthandlers.AcademicDegreeHandler,
	dh resthandlers.DatabaseHandler,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:         platform.NewRESTAPIHelper("RESTAPI"),
		academicRankHandler:   arh,
		academicDegreeHandler: adh,
		databaseHandler:       dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	academicRankHandler   resthandlers.AcademicRankHandler
	academicDegreeHandler resthandlers.AcademicDegreeHandler
	databaseHandler       resthandlers.DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/academic-ranks", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.academicRankHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.academicRankHandler.Add),
	})
	ra.NewRoute("/academic-rank/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.academicRankHandler.Find),
		http.MethodDelete: ra.NewDefaultHandler(ra.academicRankHandler.Delete),
		http.MethodPut:    ra.NewDefaultHandler(ra.academicRankHandler.Update),
	})

	ra.NewRoute("/academic-degrees", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.academicDegreeHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.academicDegreeHandler.Add),
	})
	ra.NewRoute("/academic-degree/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.academicDegreeHandler.Find),
		http.MethodDelete: ra.NewDefaultHandler(ra.academicDegreeHandler.Delete),
		http.MethodPut:    ra.NewDefaultHandler(ra.academicDegreeHandler.Update),
	})

	http.HandleFunc("/seed", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.databaseHandler.Seed)(r.Context(), w, r)
	})

	log.Printf("Server start at port %d \n", port)

	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
