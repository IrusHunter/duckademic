package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/employee/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

func NewRESTAPI(
	arh resthandlers.AcademicRankHandler,
	adh resthandlers.AcademicDegreeHandler,
	eh resthandlers.EmployeeHandler,
	th resthandlers.TeacherHandler,
	dh resthandlers.DatabaseHandler,
	jwtSecrete []byte,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:         platform.NewRESTAPIHelperWithAuth("RESTAPI", jwtSecrete),
		academicRankHandler:   arh,
		academicDegreeHandler: adh,
		employeeHandler:       eh,
		teacherHandler:        th,
		databaseHandler:       dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	academicRankHandler   resthandlers.AcademicRankHandler
	academicDegreeHandler resthandlers.AcademicDegreeHandler
	employeeHandler       resthandlers.EmployeeHandler
	teacherHandler        resthandlers.TeacherHandler
	databaseHandler       resthandlers.DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/academic-ranks", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.academicRankHandler.GetAll, []string{"employee.academic_rank"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.academicRankHandler.Add, []string{"employee.academic_rank"}),
	})
	ra.NewRoute("/academic-rank/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.academicRankHandler.Find, []string{"employee.academic_rank"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.academicRankHandler.Delete, []string{"employee.academic_rank"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.academicRankHandler.Update, []string{"employee.academic_rank"}),
	})

	ra.NewRoute("/academic-degrees", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.academicDegreeHandler.GetAll, []string{"employee.academic_degree"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.academicDegreeHandler.Add, []string{"employee.academic_degree"}),
	})
	ra.NewRoute("/academic-degree/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.academicDegreeHandler.Find, []string{"employee.academic_degree"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.academicDegreeHandler.Delete, []string{"employee.academic_degree"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.academicDegreeHandler.Update, []string{"employee.academic_degree"}),
	})

	ra.NewRoute("/employees", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.employeeHandler.GetAll, []string{"employee.employee"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.employeeHandler.Add, []string{"employee.employee"}),
	})
	ra.NewRoute("/employee/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.employeeHandler.Find, []string{"employee.employee"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.employeeHandler.Delete, []string{"employee.employee"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.employeeHandler.Update, []string{"employee.employee"}),
	})

	ra.NewRoute("/teachers", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.teacherHandler.GetAll, []string{"employee.teacher"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.teacherHandler.Add, []string{"employee.teacher"}),
	})
	ra.NewRoute("/teacher/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.teacherHandler.Find, []string{"employee.teacher"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.teacherHandler.Delete, []string{"employee.teacher"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.teacherHandler.Update, []string{"employee.teacher"}),
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
		{Name: "employee.academic_rank"},
		{Name: "employee.academic_degree"},
		{Name: "employee.employee"},
		{Name: "employee.teacher"},
	}
}
