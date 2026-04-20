package main

import (
	"context"
	"log"
	"time"

	"github.com/IrusHunter/duckademic/services/student/repositories"
	resthandlers "github.com/IrusHunter/duckademic/services/student/rest_handlers"
	"github.com/IrusHunter/duckademic/services/student/services"
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

	studentRepository := repositories.NewStudentRepository(database)
	semesterRepository := repositories.NewSemesterRepository(database)

	studentService := services.NewStudentService(studentRepository, semesterRepository, eventBus)
	semesterService := services.NewSemesterService(semesterRepository, eventBus)

	studentHandler := resthandlers.NewStudentHandler(studentService)
	semesterHandler := resthandlers.NewSemesterHandler(semesterService)
	databaseHandler := resthandlers.NewDatabaseHandler(studentService, semesterService)

	jwtSecret := envutil.GetStringFromENV("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET not specified in the .env file")
	}
	restapi := NewRESTAPI(studentHandler, semesterHandler, databaseHandler, []byte(jwtSecret))

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
