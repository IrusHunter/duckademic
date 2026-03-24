package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/curriculum/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

func NewRESTAPI(
	ch resthandlers.CurriculumHandler,
	sh resthandlers.SemesterHandler,
	dh resthandlers.DatabaseHandler,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:     platform.NewRESTAPIHelper("RESTAPI"),
		curriculumHandler: ch,
		semesterHandler:   sh,
		databaseHandler:   dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	curriculumHandler resthandlers.CurriculumHandler
	semesterHandler   resthandlers.SemesterHandler
	databaseHandler   resthandlers.DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/curriculums", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.curriculumHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.curriculumHandler.Add),
	})
	ra.NewRoute("/curriculum/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.curriculumHandler.Find),
		http.MethodDelete: ra.NewDefaultHandler(ra.curriculumHandler.Delete),
		http.MethodPut:    ra.NewDefaultHandler(ra.curriculumHandler.Update),
	})

	ra.NewRoute("/semesters", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.semesterHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.semesterHandler.Add),
	})
	ra.NewRoute("/semester/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.semesterHandler.Find),
		http.MethodDelete: ra.NewDefaultHandler(ra.semesterHandler.Delete),
		http.MethodPut:    ra.NewDefaultHandler(ra.semesterHandler.Update),
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
