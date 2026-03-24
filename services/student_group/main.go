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

	studentRepository := repositories.NewStudentRepository(database)

	studentService := services.NewStudentService(studentRepository, eventBus)

	studentHandler := resthandlers.NewStudentHandler(studentService)
	databaseHandler := resthandlers.NewDatabaseHandler(studentService)

	restapi := NewRESTAPI(studentHandler, databaseHandler)

	err = restapi.Run(port)
	log.Fatal(err)
}
