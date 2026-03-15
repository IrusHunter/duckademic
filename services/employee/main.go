package main

import (
	"log"

	"github.com/IrusHunter/duckademic/services/employees/repositories"
	resthandlers "github.com/IrusHunter/duckademic/services/employees/rest_handlers"
	"github.com/IrusHunter/duckademic/services/employees/services"
	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/IrusHunter/duckademic/shared/envutil"
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

	logger.LoadOnlyConsoleConfig()

	academicRankRepository := repositories.NewAcademicRankRepository(database)
	academicDegreeRepository := repositories.NewAcademicDegreeRepository(database)
	employeeRepository := repositories.NewEmployeeRepository(database)

	academicRankService := services.NewAcademicRankService(academicRankRepository)
	academicDegreeService := services.NewAcademicDegreeService(academicDegreeRepository)
	employeeService := services.NewEmployeeService(employeeRepository)

	academicRankHandler := resthandlers.NewAcademicRankHandler(academicRankService)
	academicDegreeHandler := resthandlers.NewAcademicDegreeHandler(academicDegreeService)
	employeeHandler := resthandlers.NewEmployeeHandler(employeeService)
	databaseHandler := resthandlers.NewDatabaseHandler(academicRankService, academicDegreeService, employeeService)

	restapi := NewRESTAPI(academicRankHandler, academicDegreeHandler, employeeHandler, databaseHandler)

	err = restapi.Run(port)
	log.Fatal(err)
}
