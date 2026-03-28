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
	lth resthandlers.LessonTypeHandler,
	disH resthandlers.DisciplineHandler,
	tlh resthandlers.TeacherLoadHandler,
	dh resthandlers.DatabaseHandler,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:      platform.NewRESTAPIHelper("RESTAPI"),
		teacherHandler:     th,
		groupCohortHandler: gch,
		lessonTypeHandler:  lth,
		disciplineHandler:  disH,
		teacherLoadHandler: tlh,
		databaseHandler:    dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	teacherHandler     resthandlers.TeacherHandler
	groupCohortHandler resthandlers.GroupCohortHandler
	databaseHandler    resthandlers.DatabaseHandler
	lessonTypeHandler  resthandlers.LessonTypeHandler
	disciplineHandler  resthandlers.DisciplineHandler
	teacherLoadHandler resthandlers.TeacherLoadHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/teachers", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.teacherHandler.GetAll),
	})

	ra.NewRoute("/group-cohorts", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.groupCohortHandler.GetAll),
	})

	ra.NewRoute("/disciplines", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.disciplineHandler.GetAll),
	})

	ra.NewRoute("/lesson-types", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.lessonTypeHandler.GetAll),
	})

	ra.NewRoute("/teacher-loads", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.teacherLoadHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.teacherLoadHandler.Add),
	})
	ra.NewRoute("/teacher-load/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.teacherLoadHandler.Find),
		http.MethodPut:    ra.NewDefaultHandler(ra.teacherLoadHandler.Update),
		http.MethodDelete: ra.NewDefaultHandler(ra.teacherLoadHandler.Delete),
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
