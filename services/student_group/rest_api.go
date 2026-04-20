package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/student_group/rest_handlers"
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
	gch resthandlers.GroupCohortHandler,
	sgh resthandlers.StudentGroupHandler,
	gmh resthandlers.GroupMemberHandler,
	lth resthandlers.LessonTypeHandler,
	disH resthandlers.DisciplineHandler,
	gcah resthandlers.GroupCohortAssignmentHandler,
	dh resthandlers.DatabaseHandler,
	jwtSecret []byte,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:                platform.NewRESTAPIHelperWithAuth("RESTAPI", jwtSecret),
		studentHandler:               sh,
		groupCohortHandler:           gch,
		semesterHandler:              semH,
		studentGroupHandler:          sgh,
		groupMembersHandler:          gmh,
		databaseHandler:              dh,
		disciplineHandler:            disH,
		lessonTypeHandler:            lth,
		groupCohortAssignmentHandler: gcah,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	studentHandler               resthandlers.StudentHandler
	semesterHandler              resthandlers.SemesterHandler
	groupCohortHandler           resthandlers.GroupCohortHandler
	studentGroupHandler          resthandlers.StudentGroupHandler
	groupMembersHandler          resthandlers.GroupMemberHandler
	disciplineHandler            resthandlers.DisciplineHandler
	lessonTypeHandler            resthandlers.LessonTypeHandler
	groupCohortAssignmentHandler resthandlers.GroupCohortAssignmentHandler
	databaseHandler              resthandlers.DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/students", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.studentHandler.GetAll, []string{"student_group.student"}),
	})

	ra.NewRoute("/semesters", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.semesterHandler.GetAll, []string{"student_group.semester"}),
	})

	ra.NewRoute("/group-cohorts", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.groupCohortHandler.GetAll, []string{"student_group.group_cohort"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.groupCohortHandler.Add, []string{"student_group.group_cohort"}),
	})
	ra.NewRoute("/group-cohort/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.groupCohortHandler.Find, []string{"student_group.group_cohort"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.groupCohortHandler.Update, []string{"student_group.group_cohort"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.groupCohortHandler.Delete, []string{"student_group.group_cohort"}),
	})

	ra.NewRoute("/student-groups", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.studentGroupHandler.GetAll, []string{"student_group.student_group"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.studentGroupHandler.Add, []string{"student_group.student_group"}),
	})
	ra.NewRoute("/student-group/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.studentGroupHandler.Find, []string{"student_group.student_group"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.studentGroupHandler.Update, []string{"student_group.student_group"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.studentGroupHandler.Delete, []string{"student_group.student_group"}),
	})

	ra.NewRoute("/group-members", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.groupMembersHandler.GetAll, []string{"student_group.group_member"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.groupMembersHandler.Add, []string{"student_group.group_member"}),
	})
	ra.NewRoute("/group-member/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.groupMembersHandler.Find, []string{"student_group.group_member"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.groupMembersHandler.Update, []string{"student_group.group_member"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.groupMembersHandler.Delete, []string{"student_group.group_member"}),
	})

	ra.NewRoute("/disciplines", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.disciplineHandler.GetAll, []string{"student_group.discipline"}),
	})

	ra.NewRoute("/lesson-types", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.lessonTypeHandler.GetAll, []string{"student_group.lesson_type"}),
	})

	ra.NewRoute("/group-cohort-assignments", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(
			ra.groupCohortAssignmentHandler.GetAll, []string{"student_group.group_cohort_assignment"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(
			ra.groupCohortAssignmentHandler.Add, []string{"student_group.group_cohort_assignment"}),
	})
	ra.NewRoute("/group-cohort-assignment/{id}", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(
			ra.groupCohortAssignmentHandler.Find, []string{"student_group.group_cohort_assignment"}),
		http.MethodPut: ra.NewDefaultHandlerWithAuth(
			ra.groupCohortAssignmentHandler.Update, []string{"student_group.group_cohort_assignment"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(
			ra.groupCohortAssignmentHandler.Delete, []string{"student_group.group_cohort_assignment"}),
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
		{Name: "student_group.student"},
		{Name: "student_group.semester"},
		{Name: "student_group.group_cohort"},
		{Name: "student_group.student_group"},
		{Name: "student_group.group_member"},
		{Name: "student_group.discipline"},
		{Name: "student_group.lesson_type"},
		{Name: "student_group.group_cohort_assignment"},
	}
}
