package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/schedule_generator/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

func NewRESTAPI(
	sgh resthandlers.ScheduleGeneratorHandler,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:            platform.NewRESTAPIHelper("RESTAPI"),
		scheduleGeneratorHandler: sgh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	scheduleGeneratorHandler resthandlers.ScheduleGeneratorHandler
}

func (ra *restapi) Run(port int) error {
	http.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.CreateGenerator)(r.Context(), w, r)
	})

	http.HandleFunc("/set-teachers", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.SetTeachers)(r.Context(), w, r)
	})

	ra.NewRoute("/default-generator-config", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.scheduleGeneratorHandler.GetDefaultConfig),
	})

	log.Printf("Server start at port %d \n", port)

	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
