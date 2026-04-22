package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/schedule/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/events"
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
	semH resthandlers.SemesterHandler,
	sdh resthandlers.SemesterDisciplineHandler,
	dh resthandlers.DatabaseHandler,
	jwtSecrete []byte,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:                platform.NewRESTAPIHelperWithAuth("RESTAPI", jwtSecrete),
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
		semesterHandler:              semH,
		semesterDisciplineHandler:    sdh,
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
	semesterHandler              resthandlers.SemesterHandler
	semesterDisciplineHandler    resthandlers.SemesterDisciplineHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/academic-ranks", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.academicRankHandler.GetAll, []string{"schedule.academic_rank"}),
	})
	ra.NewRoute("/academic-rank/{id}", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.academicRankHandler.Find, []string{"schedule.academic_rank"}),
		http.MethodPut: ra.NewDefaultHandlerWithAuth(ra.academicRankHandler.Update, []string{"schedule.academic_rank"}),
	})

	ra.NewRoute("/teachers", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.teacherHandler.GetAll, []string{"schedule.teacher"}),
	})

	ra.NewRoute("/disciplines", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.disciplineHandler.GetAll, []string{"schedule.discipline"}),
	})

	ra.NewRoute("/lesson-types", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.lessonTypeHandler.GetAll, []string{"schedule.lesson_type"}),
	})
	ra.NewRoute("/lesson-type/{id}", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.lessonTypeHandler.Find, []string{"schedule.lesson_type"}),
		http.MethodPut: ra.NewDefaultHandlerWithAuth(ra.lessonTypeHandler.Update, []string{"schedule.lesson_type"}),
	})

	ra.NewRoute("/lesson-type-assignments", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(
			ra.lessonTypeAssignmentHandler.GetAll, []string{"schedule.lesson_type_assignment"}),
	})

	ra.NewRoute("/students", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.studentHandler.GetAll, []string{"schedule.student"}),
	})

	ra.NewRoute("/student-groups", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.studentGroupHandler.GetAll, []string{"schedule.student_group"}),
	})

	ra.NewRoute("/group-members", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.groupMemberHandler.GetAll, []string{"schedule.group_member"}),
	})

	ra.NewRoute("/group-cohorts", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.groupCohortHandler.GetAll, []string{"schedule.group_cohort"}),
	})

	ra.NewRoute("/group-cohort-assignments", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(
			ra.groupCohortAssignmentHandler.GetAll, []string{"schedule.group_cohort_assignment"}),
	})

	ra.NewRoute("/teacher-loads", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.teacherLoadHandler.GetAll, []string{"schedule.teacher_load"}),
	})

	ra.NewRoute("/classrooms", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.classroomHandler.GetAll, []string{"schedule.classroom"}),
	})

	ra.NewRoute("/study-loads", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.studyLoadHandler.GetAll, []string{"schedule.study_load"}),
	})

	ra.NewRoute("/lesson-slots", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.lessonSlotHandler.GetAll, []string{"schedule.lesson_slot"}),
	})

	ra.NewRoute("/lesson-occurrences", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.lessonOccurrenceHandler.GetAll, []string{"schedule.lesson_occurrence"}),
	})

	ra.NewRoute("/semesters", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.semesterHandler.GetAll, []string{"schedule.semester"}),
	})

	ra.NewRoute("/semester-disciplines", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandlerWithAuth(ra.semesterHandler.GetAll, []string{"schedule.semester-discipline"}),
	})

	http.HandleFunc("/load-data-into-generator", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.databaseHandler.LoadDataIntoGenerator)(r.Context(), w, r)
	})
	http.HandleFunc("/extract-data-from-generator", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.databaseHandler.ExtractDataFromGenerator)(r.Context(), w, r)
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
		{Name: "schedule.academic_rank"},
		{Name: "schedule.teacher"},
		{Name: "schedule.discipline"},
		{Name: "schedule.lesson_type"},
		{Name: "schedule.lesson_type_assignment"},
		{Name: "schedule.student"},
		{Name: "schedule.student_group"},
		{Name: "schedule.group_member"},
		{Name: "schedule.group_cohort"},
		{Name: "schedule.group_cohort_assignment"},
		{Name: "schedule.teacher_load"},
		{Name: "schedule.classroom"},
		{Name: "schedule.study_load"},
		{Name: "schedule.lesson_slot"},
		{Name: "schedule.lesson_occurrence"},
		{Name: "schedule.semester"},
		{Name: "schedule.semester-discipline"},
	}
}
