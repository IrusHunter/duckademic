package main

import (
	"log"

	resthandlers "github.com/IrusHunter/duckademic/services/schedule_generator/rest_handlers"
	"github.com/IrusHunter/duckademic/services/schedule_generator/services"
	"github.com/IrusHunter/duckademic/shared/envutil"
	"github.com/IrusHunter/duckademic/shared/logger"
)

func main() {
	// testGeneration()

	if err := envutil.LoadENV(); err != nil {
		log.Fatalf(".env load failed: %s", err.Error())
	}

	port, err := envutil.GetIntFromENV("PORT")
	if err != nil {
		log.Fatalf("Can't get port value: %s", err.Error())
	}

	logger.LoadDefaultLogConfig()

	generatorConfigService := services.NewGeneratorConfigService()
	teacherService := services.NewTeacherService()

	scheduleGeneratorHandler := resthandlers.NewScheduleGeneratorHandler(generatorConfigService, teacherService)

	restapi := NewRESTAPI(scheduleGeneratorHandler)

	err = restapi.Run(port)
	log.Fatal(err)
}
