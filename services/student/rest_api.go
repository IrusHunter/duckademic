package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/student/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

func NewRESTAPI(
	sh resthandlers.StudentHandler,
	semH resthandlers.SemesterHandler,
	dh resthandlers.DatabaseHandler,
	jwtSecret []byte,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:   platform.NewRESTAPIHelperWithAuth("RESTAPI", jwtSecret),
		studentHandler:  sh,
		semesterHandler: semH,
		databaseHandler: dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	studentHandler  resthandlers.StudentHandler
	semesterHandler resthandlers.SemesterHandler
	databaseHandler resthandlers.DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/students", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.studentHandler.GetAll, []string{"student.student"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.studentHandler.Add, []string{"student.student"}),
	})
	ra.NewRoute("/student/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.studentHandler.Find, []string{"student.student"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.studentHandler.Delete, []string{"student.student"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.studentHandler.Update, []string{"student.student"}),
	})

	ra.NewRoute("/semesters", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.semesterHandler.GetAll, []string{"student.semester"}),
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

func BuildAccessPermissions() []events.AccessPermissionRE {
	return []events.AccessPermissionRE{
		{Name: "student.student"},
		{Name: "student.semester"},
	}
}
