package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/schedule/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

func NewRESTAPI(
	arh resthandlers.AcademicRankHandler,
	th resthandlers.TeacherHandler,
	disH resthandlers.DisciplineHandler,
	lth resthandlers.LessonTypeHandler,
	ltah resthandlers.LessonTypeAssignmentHandler,
	sh resthandlers.StudentHandler,
	sgh resthandlers.StudentGroupHandler,
	gmh resthandlers.GroupMemberHandler,
	tlh resthandlers.TeacherLoadHandler,
	gch resthandlers.GroupCohortHandler,
	gcah resthandlers.GroupCohortAssignmentHandler,
	ch resthandlers.ClassroomHandler,
	slh resthandlers.StudyLoadHandler,
	lsh resthandlers.LessonSlotHandler,
	loh resthandlers.LessonOccurrenceHandler,
	dh resthandlers.DatabaseHandler,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:                platform.NewRESTAPIHelper("RESTAPI"),
		academicRankHandler:          arh,
		teacherHandler:               th,
		databaseHandler:              dh,
		lessonTypeHandler:            lth,
		lessonTypeAssignmentHandler:  ltah,
		disciplineHandler:            disH,
		studentHandler:               sh,
		studentGroupHandler:          sgh,
		groupMemberHandler:           gmh,
		teacherLoadHandler:           tlh,
		groupCohortHandler:           gch,
		groupCohortAssignmentHandler: gcah,
		classroomHandler:             ch,
		studyLoadHandler:             slh,
		lessonSlotHandler:            lsh,
		lessonOccurrenceHandler:      loh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	academicRankHandler          resthandlers.AcademicRankHandler
	teacherHandler               resthandlers.TeacherHandler
	disciplineHandler            resthandlers.DisciplineHandler
	lessonTypeHandler            resthandlers.LessonTypeHandler
	lessonTypeAssignmentHandler  resthandlers.LessonTypeAssignmentHandler
	studentHandler               resthandlers.StudentHandler
	studentGroupHandler          resthandlers.StudentGroupHandler
	groupMemberHandler           resthandlers.GroupMemberHandler
	databaseHandler              resthandlers.DatabaseHandler
	groupCohortHandler           resthandlers.GroupCohortHandler
	groupCohortAssignmentHandler resthandlers.GroupCohortAssignmentHandler
	teacherLoadHandler           resthandlers.TeacherLoadHandler
	classroomHandler             resthandlers.ClassroomHandler
	studyLoadHandler             resthandlers.StudyLoadHandler
	lessonSlotHandler            resthandlers.LessonSlotHandler
	lessonOccurrenceHandler      resthandlers.LessonOccurrenceHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/academic-ranks", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.academicRankHandler.GetAll),
	})
	ra.NewRoute("/academic-rank/{id}", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.academicRankHandler.Find),
		http.MethodPut: ra.NewDefaultHandler(ra.academicRankHandler.Update),
	})

	ra.NewRoute("/teachers", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.teacherHandler.GetAll),
	})

	ra.NewRoute("/disciplines", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.disciplineHandler.GetAll),
	})

	ra.NewRoute("/lesson-types", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.lessonTypeHandler.GetAll),
	})
	ra.NewRoute("/lesson-type/{id}", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.lessonTypeHandler.Find),
		http.MethodPut: ra.NewDefaultHandler(ra.lessonTypeHandler.Update),
	})

	ra.NewRoute("/lesson-type-assignments", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.lessonTypeAssignmentHandler.GetAll),
	})

	ra.NewRoute("/students", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.studentHandler.GetAll),
	})

	ra.NewRoute("/student-groups", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.studentGroupHandler.GetAll),
	})

	ra.NewRoute("/group-members", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.groupMemberHandler.GetAll),
	})

	ra.NewRoute("/group-cohorts", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.groupCohortHandler.GetAll),
	})

	ra.NewRoute("/group-cohort-assignments", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.groupCohortAssignmentHandler.GetAll),
	})

	ra.NewRoute("/teacher-loads", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.teacherLoadHandler.GetAll),
	})

	ra.NewRoute("/classrooms", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.classroomHandler.GetAll),
	})

	http.HandleFunc("/seed", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.databaseHandler.Seed)(r.Context(), w, r)
	})
	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.databaseHandler.Clear)(r.Context(), w, r)
	})

	ra.NewRoute("/study-loads", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.studyLoadHandler.GetAll),
	})

	ra.NewRoute("/lesson-slots", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.lessonSlotHandler.GetAll),
	})

	ra.NewRoute("/lesson-occurrences", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.lessonOccurrenceHandler.GetAll),
	})

	log.Printf("Server start at port %d \n", port)

	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
