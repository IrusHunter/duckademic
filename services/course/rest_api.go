package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/course/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

func NewRESTAPI(
	sh resthandlers.StudentHandler,
	th resthandlers.TeacherHandler,
	dh resthandlers.DatabaseHandler,
	jwtSecret []byte,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:   platform.NewRESTAPIHelperWithAuth("RESTAPI", jwtSecret),
		studentHandler:  sh,
		teacherHandler:  th,
		databaseHandler: dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	studentHandler  resthandlers.StudentHandler
	teacherHandler  resthandlers.TeacherHandler
	databaseHandler resthandlers.DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/students", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.studentHandler.GetAll, []string{"course.student"}),
	})

	ra.NewRoute("/teachers", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.teacherHandler.GetAll, []string{"course.teacher"}),
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
		{Name: "course.student"},
		{Name: "course.teacher"},
	}
}
