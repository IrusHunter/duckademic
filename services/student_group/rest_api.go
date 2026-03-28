package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/student_group/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

func NewRESTAPI(
	sh resthandlers.StudentHandler,
	semH resthandlers.SemesterHandler,
	gch resthandlers.GroupCohortHandler,
	sgh resthandlers.StudentGroupHandler,
	gmh resthandlers.GroupMemberHandler,
	dh resthandlers.DatabaseHandler,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:       platform.NewRESTAPIHelper("RESTAPI"),
		studentHandler:      sh,
		groupCohortHandler:  gch,
		semesterHandler:     semH,
		studentGroupHandler: sgh,
		groupMembersHandler: gmh,
		databaseHandler:     dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	studentHandler      resthandlers.StudentHandler
	semesterHandler     resthandlers.SemesterHandler
	groupCohortHandler  resthandlers.GroupCohortHandler
	studentGroupHandler resthandlers.StudentGroupHandler
	groupMembersHandler resthandlers.GroupMemberHandler
	databaseHandler     resthandlers.DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/students", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.studentHandler.GetAll),
	})
	ra.NewRoute("/student/{id}", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.studentHandler.Find),
		http.MethodPut: ra.NewDefaultHandler(ra.studentHandler.Update),
	})

	ra.NewRoute("/semesters", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.studentHandler.GetAll),
	})
	ra.NewRoute("/semester/{id}", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.semesterHandler.Find),
		http.MethodPut: ra.NewDefaultHandler(ra.semesterHandler.Update),
	})

	ra.NewRoute("/group-cohorts", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.groupCohortHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.groupCohortHandler.Add),
	})
	ra.NewRoute("/group-cohort/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.groupCohortHandler.Find),
		http.MethodPut:    ra.NewDefaultHandler(ra.groupCohortHandler.Update),
		http.MethodDelete: ra.NewDefaultHandler(ra.groupCohortHandler.Delete),
	})

	ra.NewRoute("/student-groups", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.studentGroupHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.studentGroupHandler.Add),
	})
	ra.NewRoute("/student-group/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.studentGroupHandler.Find),
		http.MethodPut:    ra.NewDefaultHandler(ra.studentGroupHandler.Update),
		http.MethodDelete: ra.NewDefaultHandler(ra.studentGroupHandler.Delete),
	})

	ra.NewRoute("/group-members", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.groupMembersHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.groupMembersHandler.Add),
	})
	ra.NewRoute("/group-member/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.groupMembersHandler.Find),
		http.MethodPut:    ra.NewDefaultHandler(ra.groupMembersHandler.Update),
		http.MethodDelete: ra.NewDefaultHandler(ra.groupMembersHandler.Delete),
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
