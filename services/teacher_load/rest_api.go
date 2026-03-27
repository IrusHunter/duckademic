package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/teacher_load/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

func NewRESTAPI(
	th resthandlers.TeacherHandler,
	gch resthandlers.GroupCohortHandler,
	dh resthandlers.DatabaseHandler,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:      platform.NewRESTAPIHelper("RESTAPI"),
		teacherHandler:     th,
		groupCohortHandler: gch,
		databaseHandler:    dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	teacherHandler     resthandlers.TeacherHandler
	groupCohortHandler resthandlers.GroupCohortHandler
	databaseHandler    resthandlers.DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/teachers", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.teacherHandler.GetAll),
	})

	ra.NewRoute("/group_cohorts", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.groupCohortHandler.GetAll),
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
