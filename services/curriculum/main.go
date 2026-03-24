package main

import (
	"log"

	"github.com/IrusHunter/duckademic/services/curriculum/repositories"
	resthandlers "github.com/IrusHunter/duckademic/services/curriculum/rest_handlers"
	"github.com/IrusHunter/duckademic/services/curriculum/services"
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

	logger.LoadDefaultLogConfig()

	// rdc, err := events.NewDefaultRedisConnection()
	// if err != nil {
	// 	log.Fatalf("Can't connect to redis: %v", err)
	// }
	// eventBus := events.NewEventBus(rdc)

	curriculumRepository := repositories.NewCurriculumRepository(database)

	curriculumService := services.NewCurriculumService(curriculumRepository)

	curriculumHandler := resthandlers.NewCurriculumHandler(curriculumService)
	databaseHandler := resthandlers.NewDatabaseHandler(curriculumService)

	restapi := NewRESTAPI(curriculumHandler, databaseHandler)

	err = restapi.Run(port)
	log.Fatal(err)
}
