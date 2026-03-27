package main

import (
	"log"

	"github.com/IrusHunter/duckademic/services/schedule/repositories"
	resthandlers "github.com/IrusHunter/duckademic/services/schedule/rest_handlers"
	"github.com/IrusHunter/duckademic/services/schedule/services"
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

	academicRankRepository := repositories.NewAcademicRankRepository(database)
	teacherRepository := repositories.NewTeacherRepository(database)
	disciplineRepository := repositories.NewDisciplineRepository(database)
	lessonTypeRepository := repositories.NewLessonTypeRepository(database)
	lessonTypeAssignmentRepository := repositories.NewLessonTypeAssignmentRepository(database)
	studentRepository := repositories.NewStudentRepository(database)
	studentGroupRepository := repositories.NewStudentGroupRepository(database)
	groupMemberRepository := repositories.NewGroupMemberRepository(database)

	academicRankService := services.NewAcademicRankService(academicRankRepository, eventBus)
	teacherService := services.NewTeacherService(teacherRepository, eventBus)
	disciplineService := services.NewDisciplineService(disciplineRepository, eventBus)
	lessonTypeService := services.NewLessonTypeService(lessonTypeRepository, eventBus)
	lessonTypeAssignmentService := services.NewLessonTypeAssignmentService(lessonTypeAssignmentRepository,
		lessonTypeRepository, disciplineRepository, eventBus)
	studentService := services.NewStudentService(studentRepository, eventBus)
	studentGroupService := services.NewStudentGroupService(studentGroupRepository, eventBus)
	groupMemberService := services.NewGroupMemberService(groupMemberRepository, eventBus)

	academicRankHandler := resthandlers.NewAcademicRankHandler(academicRankService)
	teacherHandler := resthandlers.NewTeacherHandler(teacherService)
	disciplineHandler := resthandlers.NewDisciplineHandler(disciplineService)
	lessonTypeHandler := resthandlers.NewLessonTypeHandler(lessonTypeService)
	lessonTypeAssignmentHandler := resthandlers.NewLessonTypeAssignmentHandler(lessonTypeAssignmentService)
	databaseHandler := resthandlers.NewDatabaseHandler(academicRankService)
	studentHandler := resthandlers.NewStudentHandler(studentService)
	studentGroupHandler := resthandlers.NewStudentGroupHandler(studentGroupService)
	groupMemberHandler := resthandlers.NewGroupMemberHandler(groupMemberService)

	restapi := NewRESTAPI(academicRankHandler, teacherHandler, disciplineHandler, lessonTypeHandler,
		lessonTypeAssignmentHandler, studentHandler, studentGroupHandler, groupMemberHandler, databaseHandler)

	err = restapi.Run(port)
	log.Fatal(err)
}
