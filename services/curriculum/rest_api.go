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
	lth resthandlers.LessonTypeHandler,
	disH resthandlers.DisciplineHandler,
	ltah resthandlers.LessonTypeAssignmentHandler,
	sdh resthandlers.SemesterDisciplineHandler,
	dh resthandlers.DatabaseHandler,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:               platform.NewRESTAPIHelper("RESTAPI"),
		curriculumHandler:           ch,
		semesterHandler:             sh,
		lessonTypeHandler:           lth,
		disciplineHandler:           disH,
		lessonTypeAssignmentHandler: ltah,
		semesterDisciplineHandler:   sdh,
		databaseHandler:             dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	curriculumHandler           resthandlers.CurriculumHandler
	semesterHandler             resthandlers.SemesterHandler
	lessonTypeHandler           resthandlers.LessonTypeHandler
	disciplineHandler           resthandlers.DisciplineHandler
	lessonTypeAssignmentHandler resthandlers.LessonTypeAssignmentHandler
	semesterDisciplineHandler   resthandlers.SemesterDisciplineHandler
	databaseHandler             resthandlers.DatabaseHandler
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

	ra.NewRoute("/lesson-types", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.lessonTypeHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.lessonTypeHandler.Add),
	})
	ra.NewRoute("/lesson-type/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.lessonTypeHandler.Find),
		http.MethodDelete: ra.NewDefaultHandler(ra.lessonTypeHandler.Delete),
		http.MethodPut:    ra.NewDefaultHandler(ra.lessonTypeHandler.Update),
	})

	ra.NewRoute("/disciplines", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.disciplineHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.disciplineHandler.Add),
	})
	ra.NewRoute("/discipline/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.disciplineHandler.Find),
		http.MethodDelete: ra.NewDefaultHandler(ra.disciplineHandler.Delete),
		http.MethodPut:    ra.NewDefaultHandler(ra.disciplineHandler.Update),
	})

	ra.NewRoute("/lesson-type-assignments", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.lessonTypeAssignmentHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.lessonTypeAssignmentHandler.Add),
	})
	ra.NewRoute("/lesson-type-assignment/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.lessonTypeAssignmentHandler.Find),
		http.MethodDelete: ra.NewDefaultHandler(ra.lessonTypeAssignmentHandler.Delete),
		http.MethodPut:    ra.NewDefaultHandler(ra.lessonTypeAssignmentHandler.Update),
	})

	ra.NewRoute("/semester-disciplines", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.semesterDisciplineHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.semesterDisciplineHandler.Add),
	})
	ra.NewRoute("/semester-discipline/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.semesterDisciplineHandler.Find),
		http.MethodDelete: ra.NewDefaultHandler(ra.semesterDisciplineHandler.Delete),
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
