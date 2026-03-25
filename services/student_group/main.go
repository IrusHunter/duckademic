package main

import (
	"log"

	"github.com/IrusHunter/duckademic/services/student_group/repositories"
	resthandlers "github.com/IrusHunter/duckademic/services/student_group/rest_handlers"
	"github.com/IrusHunter/duckademic/services/student_group/services"
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

	semesterRepository := repositories.NewSemesterRepository(database)
	studentRepository := repositories.NewStudentRepository(database)
	groupCohortRepository := repositories.NewGroupCohortRepository(database)
	studentGroupRepository := repositories.NewStudentGroupRepository(database)

	semesterService := services.NewSemesterService(semesterRepository, eventBus)
	studentService := services.NewStudentService(studentRepository, eventBus)
	groupCohortService := services.NewGroupCohortService(groupCohortRepository, semesterRepository)
	studentGroupService := services.NewStudentGroupService(studentGroupRepository, groupCohortRepository)

	semesterHandler := resthandlers.NewSemesterHandler(semesterService)
	studentHandler := resthandlers.NewStudentHandler(studentService)
	groupCohortHandler := resthandlers.NewGroupCohortHandler(groupCohortService)
	studentGroupHandler := resthandlers.NewStudentGroupHandler(studentGroupService)
	databaseHandler := resthandlers.NewDatabaseHandler(studentService, semesterService, groupCohortService,
		studentGroupService)

	restapi := NewRESTAPI(studentHandler, semesterHandler, groupCohortHandler, studentGroupHandler,
		databaseHandler)

	err = restapi.Run(port)
	log.Fatal(err)
}
