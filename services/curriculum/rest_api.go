package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/curriculum/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/events"
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
	jwtSecret []byte,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:               platform.NewRESTAPIHelperWithAuth("RESTAPI", jwtSecret),
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
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.curriculumHandler.GetAll, []string{"curriculum.curriculum"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.curriculumHandler.Add, []string{"curriculum.curriculum"}),
	})
	ra.NewRoute("/curriculum/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.curriculumHandler.Find, []string{"curriculum.curriculum"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.curriculumHandler.Delete, []string{"curriculum.curriculum"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.curriculumHandler.Update, []string{"curriculum.curriculum"}),
	})

	ra.NewRoute("/semesters", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.semesterHandler.GetAll, []string{"curriculum.semester"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.semesterHandler.Add, []string{"curriculum.semester"}),
	})
	ra.NewRoute("/semester/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.semesterHandler.Find, []string{"curriculum.semester"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.semesterHandler.Delete, []string{"curriculum.semester"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.semesterHandler.Update, []string{"curriculum.semester"}),
	})

	ra.NewRoute("/lesson-types", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.lessonTypeHandler.GetAll, []string{"curriculum.lesson_type"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.lessonTypeHandler.Add, []string{"curriculum.lesson_type"}),
	})
	ra.NewRoute("/lesson-type/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.lessonTypeHandler.Find, []string{"curriculum.lesson_type"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.lessonTypeHandler.Delete, []string{"curriculum.lesson_type"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.lessonTypeHandler.Update, []string{"curriculum.lesson_type"}),
	})

	ra.NewRoute("/disciplines", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.disciplineHandler.GetAll, []string{"curriculum.discipline"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.disciplineHandler.Add, []string{"curriculum.discipline"}),
	})
	ra.NewRoute("/discipline/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.disciplineHandler.Find, []string{"curriculum.discipline"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.disciplineHandler.Delete, []string{"curriculum.discipline"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.disciplineHandler.Update, []string{"curriculum.discipline"}),
	})

	ra.NewRoute("/lesson-type-assignments", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(
			ra.lessonTypeAssignmentHandler.GetAll, []string{"curriculum.lesson_type_assignment"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(
			ra.lessonTypeAssignmentHandler.Add, []string{"curriculum.lesson_type_assignment"}),
	})
	ra.NewRoute("/lesson-type-assignment/{id}", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(
			ra.lessonTypeAssignmentHandler.Find, []string{"curriculum.lesson_type_assignment"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(
			ra.lessonTypeAssignmentHandler.Delete, []string{"curriculum.lesson_type_assignment"}),
		http.MethodPut: ra.NewDefaultHandlerWithAuth(
			ra.lessonTypeAssignmentHandler.Update, []string{"curriculum.lesson_type_assignment"}),
	})

	ra.NewRoute("/semester-disciplines", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(
			ra.semesterDisciplineHandler.GetAll, []string{"curriculum.semester_discipline"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(
			ra.semesterDisciplineHandler.Add, []string{"curriculum.semester_discipline"}),
	})
	ra.NewRoute("/semester-discipline/{id}", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(
			ra.semesterDisciplineHandler.Find, []string{"curriculum.semester_discipline"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(
			ra.semesterDisciplineHandler.Delete, []string{"curriculum.semester_discipline"}),
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
		{Name: "curriculum.curriculum"},
		{Name: "curriculum.semester"},
		{Name: "curriculum.lesson_type"},
		{Name: "curriculum.discipline"},
		{Name: "curriculum.lesson_type_assignment"},
		{Name: "curriculum.semester_discipline"},
	}
}
