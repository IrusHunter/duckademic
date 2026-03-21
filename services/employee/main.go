package main

import (
	"log"

	"github.com/IrusHunter/duckademic/services/employee/repositories"
	resthandlers "github.com/IrusHunter/duckademic/services/employee/rest_handlers"
	"github.com/IrusHunter/duckademic/services/employee/services"
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

	rdc, err := events.NewDefaultRedisConnection()
	if err != nil {
		log.Fatalf("Can't connect to redis: %v", err)
	}
	eventBus := events.NewEventBus(rdc)

	logger.LoadDefaultLogConfig()

	academicRankRepository := repositories.NewAcademicRankRepository(database)
	academicDegreeRepository := repositories.NewAcademicDegreeRepository(database)
	employeeRepository := repositories.NewEmployeeRepository(database)
	teacherRepository := repositories.NewTeacherRepository(database)

	academicRankService := services.NewAcademicRankService(academicRankRepository, eventBus)
	academicDegreeService := services.NewAcademicDegreeService(academicDegreeRepository)
	employeeService := services.NewEmployeeService(employeeRepository)
	teacherService := services.NewTeacherService(teacherRepository, academicRankRepository,
		academicDegreeRepository, employeeRepository)

	academicRankHandler := resthandlers.NewAcademicRankHandler(academicRankService)
	academicDegreeHandler := resthandlers.NewAcademicDegreeHandler(academicDegreeService)
	employeeHandler := resthandlers.NewEmployeeHandler(employeeService)
	teacherHandler := resthandlers.NewTeacherHandler(teacherService)
	databaseHandler := resthandlers.NewDatabaseHandler(academicRankService, academicDegreeService,
		employeeService, teacherService)

	restapi := NewRESTAPI(academicRankHandler, academicDegreeHandler, employeeHandler,
		teacherHandler, databaseHandler)

	err = restapi.Run(port)
	log.Fatal(err)
}
