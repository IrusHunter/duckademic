package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/IrusHunter/duckademic/services/schedule/repositories"
	resthandlers "github.com/IrusHunter/duckademic/services/schedule/rest_handlers"
	"github.com/IrusHunter/duckademic/services/schedule/services"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/IrusHunter/duckademic/shared/envutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/logger"
)

func main() {
	if err := envutil.LoadENV(); err != nil {
		log.Fatalf(".env load failed: %s", err.Error())
	}

	port, err := envutil.GetIntFromENV("PORT")
	if err != nil {
		log.Fatalf("Can't get port value: %s", err.Error())
	}

	database, err := db.NewDefaultDBConnection()
	if err != nil {
		log.Fatalf("Can't connect to database: %v", err)
	}

	err = Migrate(database)
	if err != nil {
		log.Fatalf("Can't migrate the database: %s", err.Error())
	}

	logger.LoadDefaultLogConfig()

	rdc, err := events.NewDefaultRedisConnection()
	if err != nil {
		log.Fatalf("Can't connect to redis: %v", err)
	}
	eventBus := events.NewEventBus(rdc)

	scheduleGeneratorDomain := envutil.GetStringFromENV("SCHEDULE_GENERATOR_SERVICE")
	if scheduleGeneratorDomain == "" {
		log.Fatalf("SCHEDULE_GENERATOR_SERVICE not specified at .env file")
	}

	academicRankRepository := repositories.NewAcademicRankRepository(database)
	teacherRepository := repositories.NewTeacherRepository(database)
	disciplineRepository := repositories.NewDisciplineRepository(database)
	lessonTypeRepository := repositories.NewLessonTypeRepository(database)
	lessonTypeAssignmentRepository := repositories.NewLessonTypeAssignmentRepository(database)
	studentRepository := repositories.NewStudentRepository(database)
	studentGroupRepository := repositories.NewStudentGroupRepository(database)
	groupMemberRepository := repositories.NewGroupMemberRepository(database)
	teacherLoadRepository := repositories.NewTeacherLoadRepository(database)
	groupCohortRepository := repositories.NewGroupCohortRepository(database)
	groupCohortAssignmentRepository := repositories.NewGroupCohortAssignmentRepository(database)
	classroomRepository := repositories.NewClassroomRepository(database)
	studyLoadRepository := repositories.NewStudyLoadRepository(database)
	lessonSlotRepository := repositories.NewLessonSlotRepository(database)
	lessonOccurrenceRepository := repositories.NewLessonOccurrenceRepository(database)

	academicRankService := services.NewAcademicRankService(academicRankRepository, eventBus)
	teacherService := services.NewTeacherService(teacherRepository, eventBus)
	disciplineService := services.NewDisciplineService(disciplineRepository, eventBus)
	lessonTypeService := services.NewLessonTypeService(lessonTypeRepository, eventBus)
	lessonTypeAssignmentService := services.NewLessonTypeAssignmentService(lessonTypeAssignmentRepository,
		lessonTypeRepository, disciplineRepository, eventBus)
	studentService := services.NewStudentService(studentRepository, eventBus)
	studentGroupService := services.NewStudentGroupService(studentGroupRepository, eventBus)
	groupMemberService := services.NewGroupMemberService(groupMemberRepository, eventBus)
	teacherLoadService := services.NewTeacherLoadService(teacherLoadRepository, eventBus)
	groupCohortService := services.NewGroupCohortService(groupCohortRepository, eventBus)
	groupCohortAssignmentService := services.NewGroupCohortAssignmentService(groupCohortAssignmentRepository, eventBus)
	classroomService := services.NewClassroomService(classroomRepository, eventBus)
	studyLoadService := services.NewStudyLoadService(studyLoadRepository)
	lessonSlotService := services.NewLessonSlotService(lessonSlotRepository)
	lessonOccurrenceService := services.NewLessonOccurrenceService(lessonOccurrenceRepository, lessonSlotRepository)

	academicRankHandler := resthandlers.NewAcademicRankHandler(academicRankService)
	teacherHandler := resthandlers.NewTeacherHandler(teacherService)
	disciplineHandler := resthandlers.NewDisciplineHandler(disciplineService)
	lessonTypeHandler := resthandlers.NewLessonTypeHandler(lessonTypeService)
	lessonTypeAssignmentHandler := resthandlers.NewLessonTypeAssignmentHandler(lessonTypeAssignmentService)
	studentHandler := resthandlers.NewStudentHandler(studentService)
	studentGroupHandler := resthandlers.NewStudentGroupHandler(studentGroupService)
	groupMemberHandler := resthandlers.NewGroupMemberHandler(groupMemberService)
	teacherLoadHandler := resthandlers.NewTeacherLoadHandler(teacherLoadService)
	groupCohortHandler := resthandlers.NewGroupCohortHandler(groupCohortService)
	groupCohortAssignmentHandler := resthandlers.NewGroupCohortAssignmentHandler(groupCohortAssignmentService)
	classroomHandler := resthandlers.NewClassroomHandler(classroomService)
	studyLoadHandler := resthandlers.NewStudyLoadHandler(studyLoadService)
	lessonSlotHandler := resthandlers.NewLessonSlotHandler(lessonSlotService)
	lessonOccurrenceHandler := resthandlers.NewLessonOccurrenceHandler(lessonOccurrenceService)
	databaseHandler := resthandlers.NewDatabaseHandler(http.DefaultClient, scheduleGeneratorDomain, academicRankService,
		teacherService, disciplineService, lessonTypeService, lessonTypeAssignmentService, studentService, studentGroupService,
		groupMemberService, teacherLoadService, groupCohortService, groupCohortAssignmentService, classroomService,
		studyLoadService, lessonSlotService, lessonOccurrenceService)

	jwtSecret := envutil.GetStringFromENV("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET not specified in the .env file")
	}
	restapi := NewRESTAPI(academicRankHandler, teacherHandler, disciplineHandler, lessonTypeHandler,
		lessonTypeAssignmentHandler, studentHandler, studentGroupHandler, groupMemberHandler, teacherLoadHandler,
		groupCohortHandler, groupCohortAssignmentHandler, classroomHandler, studyLoadHandler, lessonSlotHandler,
		lessonOccurrenceHandler, databaseHandler, []byte(jwtSecret))

	go func() {
		time.Sleep(events.ExternalSeedCooldown)
		ctx := contextutil.SetTraceID(context.Background())
		err := eventBus.PublishAccessPermissions(ctx, BuildAccessPermissions())
		if err != nil {
			log.Fatalf("Can't publish access permissions: %s", err)
		}
	}()
	err = restapi.Run(port)
	log.Fatal(err)
}
