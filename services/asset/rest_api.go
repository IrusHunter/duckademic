package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/asset/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

func NewRESTAPI(
	ch resthandlers.ClassroomHandler,
	dh resthandlers.DatabaseHandler,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:    platform.NewRESTAPIHelper("RESTAPI"),
		classroomHandler: ch,
		databaseHandler:  dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	classroomHandler resthandlers.ClassroomHandler
	databaseHandler  resthandlers.DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/classrooms", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.classroomHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.classroomHandler.Add),
	})
	ra.NewRoute("/classroom/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.classroomHandler.Find),
		http.MethodDelete: ra.NewDefaultHandler(ra.classroomHandler.Delete),
		http.MethodPut:    ra.NewDefaultHandler(ra.classroomHandler.Update),
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
