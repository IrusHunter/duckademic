package main

import (
	"context"
	"log"
	"time"

	"github.com/IrusHunter/duckademic/services/employee/repositories"
	resthandlers "github.com/IrusHunter/duckademic/services/employee/rest_handlers"
	"github.com/IrusHunter/duckademic/services/employee/services"
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
		academicDegreeRepository, employeeRepository, eventBus)

	academicRankHandler := resthandlers.NewAcademicRankHandler(academicRankService)
	academicDegreeHandler := resthandlers.NewAcademicDegreeHandler(academicDegreeService)
	employeeHandler := resthandlers.NewEmployeeHandler(employeeService)
	teacherHandler := resthandlers.NewTeacherHandler(teacherService)
	databaseHandler := resthandlers.NewDatabaseHandler(academicRankService, academicDegreeService,
		employeeService, teacherService)

	jwtSecret := envutil.GetStringFromENV("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET not specified in the .env file")
	}
	restapi := NewRESTAPI(academicRankHandler, academicDegreeHandler, employeeHandler,
		teacherHandler, databaseHandler, []byte(jwtSecret))

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
