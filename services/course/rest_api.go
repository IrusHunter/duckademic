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
	tch resthandlers.TeacherCourseHandler,
	taskH resthandlers.TaskHandler,
	tsh resthandlers.TaskStudentHandler,
	dh resthandlers.DatabaseHandler,
	jwtSecret []byte,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:        platform.NewRESTAPIHelperWithAuth("RESTAPI", jwtSecret),
		studentHandler:       sh,
		teacherHandler:       th,
		courseHandler:        ch,
		studentCourseHandler: sch,
		teacherCourseHandler: tch,
		taskHandler:          taskH,
		taskStudentHandler:   tsh,
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
	taskHandler          resthandlers.TaskHandler
	taskStudentHandler   resthandlers.TaskStudentHandler
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

	ra.NewRoute("/tasks", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.taskHandler.GetAll, []string{"course.task"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.taskHandler.Add, []string{"course.task"}),
	})
	ra.NewRoute("/task/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.taskHandler.Find, []string{"course.task"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.taskHandler.Delete, []string{"course.task"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.taskHandler.Update, []string{"course.task"}),
	})

	ra.NewRoute("/task-students", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.taskStudentHandler.GetAll, []string{"course.task_student"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.taskStudentHandler.Add, []string{"course.task_student"}),
	})
	ra.NewRoute("/task-student/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.taskStudentHandler.Find, []string{"course.task_student"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.taskStudentHandler.Delete, []string{"course.task_student"}),
	})

	ra.NewRoute("/get-upcoming-events", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.taskStudentHandler.GetUpcomingEvents, []string{}),
	})
	ra.NewRoute("/get-courses-progress", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.studentCourseHandler.GetCourseProgress, []string{}),
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
		{Name: "course.task"},
		{Name: "course.task_student"},
	}
}
