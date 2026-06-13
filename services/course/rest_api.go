package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	resthandlers "github.com/IrusHunter/duckademic/services/course/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
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
	ra.NewRoute("/course/{id}/tasks", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.taskHandler.GetByCourseID, []string{}),
	})
	ra.NewRoute("/course/{id}/student-tasks", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.taskStudentHandler.GetForStudentInCourse, []string{}),
	})

	ra.NewRoute("/courses/student", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.courseHandler.GetStudentCoursePage, []string{}),
	})

	ra.NewRoute("/courses/teacher", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.courseHandler.GetTeacherCoursePage, []string{}),
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
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.taskHandler.GetAll, []string{}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.taskHandler.Add, []string{}),
	})
	ra.NewRoute("/task/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.taskHandler.Find, []string{}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.taskHandler.Delete, []string{}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.taskHandler.Update, []string{}),
	})
	ra.NewRoute("/task/{id}/submissions", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.taskStudentHandler.GetByTaskID, []string{}),
	})

	ra.NewRoute("/task-students", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.taskStudentHandler.GetAll, []string{}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.taskStudentHandler.Add, []string{}),
	})
	ra.NewRoute("/task-student/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.taskStudentHandler.Find, []string{}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.taskStudentHandler.Delete, []string{}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.taskStudentHandler.Update, []string{}),
	})

	ra.NewRoute("/get-upcoming-events", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.taskStudentHandler.GetUpcomingEvents, []string{}),
	})
	ra.NewRoute("/get-courses-progress", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.studentCourseHandler.GetCourseProgress, []string{}),
	})

	ra.NewRoute("/upload", map[string]platform.HandlerFunc{
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.uploadFile, []string{}),
	})

	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("data/uploads"))))

	http.HandleFunc("/seed", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.databaseHandler.Seed)(r.Context(), w, r)
	})
	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.databaseHandler.Clear)(r.Context(), w, r)
	})

	log.Printf("Server start at port %d \n", port)

	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

const maxUploadSize = 10 << 20 // 10 MB

var allowedExtensions = map[string]bool{
	".pdf": true, ".doc": true, ".docx": true, ".txt": true,
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
	".zip": true, ".pptx": true, ".xlsx": true,
}

func (ra *restapi) uploadFile(_ context.Context, w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		jsonutil.ResponseWithError(w, http.StatusBadRequest, fmt.Errorf("file too large (max 10 MB)"))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		jsonutil.ResponseWithError(w, http.StatusBadRequest, fmt.Errorf("file field is required"))
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExtensions[ext] {
		jsonutil.ResponseWithError(w, http.StatusBadRequest, fmt.Errorf("file type %q is not allowed", ext))
		return
	}

	if err := os.MkdirAll("data/uploads", 0755); err != nil {
		jsonutil.ResponseWithError(w, http.StatusInternalServerError, fmt.Errorf("storage error"))
		return
	}

	filename := uuid.New().String() + ext
	dst, err := os.Create(filepath.Join("data/uploads", filename))
	if err != nil {
		jsonutil.ResponseWithError(w, http.StatusInternalServerError, fmt.Errorf("could not save file"))
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		jsonutil.ResponseWithError(w, http.StatusInternalServerError, fmt.Errorf("could not write file"))
		return
	}

	jsonutil.ResponseWithJSON(w, http.StatusOK, map[string]string{
		"url": "/uploads/" + filename,
	})
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
