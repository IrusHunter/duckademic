package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/schedule_generator/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

func NewRESTAPI(
	sgh resthandlers.ScheduleGeneratorHandler,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:            platform.NewRESTAPIHelper("RESTAPI"),
		scheduleGeneratorHandler: sgh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	scheduleGeneratorHandler resthandlers.ScheduleGeneratorHandler
}

func (ra *restapi) Run(port int) error {
	http.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.CreateGenerator)(r.Context(), w, r)
	})

	http.HandleFunc("/set-teachers", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.SetTeachers)(r.Context(), w, r)
	})
	http.HandleFunc("/set-disciplines", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.SetDisciplines)(r.Context(), w, r)
	})
	http.HandleFunc("/set-lesson-types", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.SetLessonTypes)(r.Context(), w, r)
	})
	http.HandleFunc("/set-lesson-type-assignments", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.SetLessonTypeAssignments)(r.Context(), w, r)
	})
	http.HandleFunc("/set-student-groups", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.SetStudentGroups)(r.Context(), w, r)
	})
	http.HandleFunc("/set-study-loads", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.SetStudyLoads)(r.Context(), w, r)
	})
	http.HandleFunc("/set-classrooms", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.SetClassrooms)(r.Context(), w, r)
	})
	http.HandleFunc("/submit-and-go", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.SubmitAndGoToTheNextStep)(r.Context(), w, r)
	})
	http.HandleFunc("/generate-days-for-lesson-types", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.SetDaysForLessonTypes)(r.Context(), w, r)
	})
	http.HandleFunc("/generate-bone-lessons", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.GenerateBoneLessons)(r.Context(), w, r)
	})
	http.HandleFunc("/assign-classrooms-to-bone-lessons", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.AssignClassroomsToBoneLessons)(r.Context(), w, r)
	})
	http.HandleFunc("/build-schedule-skeleton", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.BuildScheduleSkeleton)(r.Context(), w, r)
	})
	http.HandleFunc("/add-floating-lessons", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.AddFloatingLessons)(r.Context(), w, r)
	})
	http.HandleFunc("/assign-classrooms-to-floating-lessons", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.AssignClassroomsToFloatingLessons)(r.Context(), w, r)
	})
	http.HandleFunc("/get-study-loads", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.GetStudyLoads)(r.Context(), w, r)
	})
	http.HandleFunc("/get-lessons", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.scheduleGeneratorHandler.GetLessons)(r.Context(), w, r)
	})

	ra.NewRoute("/default-generator-config", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.scheduleGeneratorHandler.GetDefaultConfig),
	})

	log.Printf("Server start at port %d \n", port)

	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
