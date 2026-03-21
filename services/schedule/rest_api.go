package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/schedule/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

func NewRESTAPI(
	arh resthandlers.AcademicRankHandler,
	th resthandlers.TeacherHandler,
	dh resthandlers.DatabaseHandler,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:       platform.NewRESTAPIHelper("RESTAPI"),
		academicRankHandler: arh,
		teacherHandler:      th,
		databaseHandler:     dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	academicRankHandler resthandlers.AcademicRankHandler
	teacherHandler      resthandlers.TeacherHandler
	databaseHandler     resthandlers.DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/academic-ranks", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.academicRankHandler.GetAll),
	})
	ra.NewRoute("/academic-rank/{id}", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.academicRankHandler.Find),
		http.MethodPut: ra.NewDefaultHandler(ra.academicRankHandler.Update),
	})
	ra.NewRoute("/teachers", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.teacherHandler.GetAll),
	})
	ra.NewRoute("/teacher/{id}", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.teacherHandler.Find),
		http.MethodPut: ra.NewDefaultHandler(ra.teacherHandler.Update),
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
