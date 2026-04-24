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
	ch resthandlers.CourseHandler,
	sch resthandlers.StudentCourseHandler,
	tsh resthandlers.TeacherCourseHandler,
	dh resthandlers.DatabaseHandler,
	jwtSecret []byte,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:        platform.NewRESTAPIHelperWithAuth("RESTAPI", jwtSecret),
		studentHandler:       sh,
		teacherHandler:       th,
		courseHandler:        ch,
		studentCourseHandler: sch,
		teacherCourseHandler: tsh,
		databaseHandler:      dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	studentHandler       resthandlers.StudentHandler
	teacherHandler       resthandlers.TeacherHandler
	courseHandler        resthandlers.CourseHandler
	studentCourseHandler resthandlers.StudentCourseHandler
	teacherCourseHandler resthandlers.TeacherCourseHandler
	databaseHandler      resthandlers.DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/students", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.studentHandler.GetAll, []string{"course.student"}),
	})

	ra.NewRoute("/teachers", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.teacherHandler.GetAll, []string{"course.teacher"}),
	})

	ra.NewRoute("/courses", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.courseHandler.GetAll, []string{"course.course"}),
	})
	ra.NewRoute("/course/{id}", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.courseHandler.Find, []string{"course.course"}),
		http.MethodPut: ra.NewDefaultHandlerWithAuth(ra.courseHandler.Update, []string{"course.course"}),
	})

	ra.NewRoute("/student-courses", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.studentCourseHandler.GetAll, []string{"course.student_course"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.studentCourseHandler.Add, []string{"course.student_course"}),
	})
	ra.NewRoute("/student-course/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.studentCourseHandler.Find, []string{"course.student_course"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.studentCourseHandler.Delete, []string{"course.student_course"}),
	})

	ra.NewRoute("/teacher-courses", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.teacherCourseHandler.GetAll, []string{"course.teacher_course"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.teacherCourseHandler.Add, []string{"course.teacher_course"}),
	})
	ra.NewRoute("/teacher-course/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.teacherCourseHandler.Find, []string{"course.teacher_course"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.teacherCourseHandler.Delete, []string{"course.teacher_course"}),
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
		{Name: "course.course"},
		{Name: "course.student_course"},
		{Name: "course.teacher_course"},
	}
}
