package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/teacher_load/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

func NewRESTAPI(
	th resthandlers.TeacherHandler,
	lth resthandlers.LessonTypeHandler,
	disH resthandlers.DisciplineHandler,
	tlh resthandlers.TeacherLoadHandler,
	dh resthandlers.DatabaseHandler,
	jwtSecret []byte,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:      platform.NewRESTAPIHelperWithAuth("RESTAPI", jwtSecret),
		teacherHandler:     th,
		lessonTypeHandler:  lth,
		disciplineHandler:  disH,
		teacherLoadHandler: tlh,
		databaseHandler:    dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	teacherHandler     resthandlers.TeacherHandler
	databaseHandler    resthandlers.DatabaseHandler
	lessonTypeHandler  resthandlers.LessonTypeHandler
	disciplineHandler  resthandlers.DisciplineHandler
	teacherLoadHandler resthandlers.TeacherLoadHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/teachers", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.teacherHandler.GetAll, []string{"teacher_load.teacher"}),
	})

	ra.NewRoute("/disciplines", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.disciplineHandler.GetAll, []string{"teacher_load.discipline"}),
	})

	ra.NewRoute("/lesson-types", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.lessonTypeHandler.GetAll, []string{"teacher_load.lesson_type"}),
	})

	ra.NewRoute("/teacher-loads", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.teacherLoadHandler.GetAll, []string{"teacher_load.teacher_load"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.teacherLoadHandler.Add, []string{"teacher_load.teacher_load"}),
	})
	ra.NewRoute("/teacher-load/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.teacherLoadHandler.Find, []string{"teacher_load.teacher_load"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.teacherLoadHandler.Update, []string{"teacher_load.teacher_load"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.teacherLoadHandler.Delete, []string{"teacher_load.teacher_load"}),
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
		{Name: "teacher_load.teacher"},
		{Name: "teacher_load.discipline"},
		{Name: "teacher_load.lesson_type"},
		{Name: "teacher_load.teacher_load"},
	}
}
