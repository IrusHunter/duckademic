package main

import (
	"context"
	"log"
	"time"

	"github.com/IrusHunter/duckademic/services/curriculum/repositories"
	resthandlers "github.com/IrusHunter/duckademic/services/curriculum/rest_handlers"
	"github.com/IrusHunter/duckademic/services/curriculum/services"
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

	curriculumRepository := repositories.NewCurriculumRepository(database)
	semesterRepository := repositories.NewSemesterRepository(database)
	lessonTypeRepository := repositories.NewLessonTypeRepository(database)
	disciplineRepository := repositories.NewDisciplineRepository(database)
	lessonTypeAssignmentRepository := repositories.NewLessonTypeAssignmentRepository(database)
	semesterDisciplineRepository := repositories.NewSemesterDisciplineRepository(database)

	curriculumService := services.NewCurriculumService(curriculumRepository)
	semesterService := services.NewSemesterService(semesterRepository, curriculumRepository, eventBus)
	lessonTypeService := services.NewLessonTypeService(lessonTypeRepository, eventBus)
	disciplineService := services.NewDisciplineService(disciplineRepository, eventBus)
	lessonTypeAssignmentService := services.NewLessonTypeAssignmentService(lessonTypeAssignmentRepository,
		lessonTypeRepository, disciplineRepository, eventBus)
	semesterDisciplineService := services.NewSemesterDisciplineService(semesterDisciplineRepository, semesterRepository,
		disciplineRepository, curriculumRepository, eventBus)

	curriculumHandler := resthandlers.NewCurriculumHandler(curriculumService)
	semesterHandler := resthandlers.NewSemesterHandler(semesterService)
	lessonTypeHandler := resthandlers.NewLessonTypeHandler(lessonTypeService)
	disciplineHandler := resthandlers.NewDisciplineHandler(disciplineService)
	lessonTypeAssignmentHandler := resthandlers.NewLessonTypeAssignmentHandler(lessonTypeAssignmentService)
	semesterDisciplineHandler := resthandlers.NewSemesterDisciplineHandler(semesterDisciplineService)
	databaseHandler := resthandlers.NewDatabaseHandler(curriculumService, semesterService, lessonTypeService,
		disciplineService, lessonTypeAssignmentService, semesterDisciplineService)

	jwtSecret := envutil.GetStringFromENV("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET not specified in the .env file")
	}
	restapi := NewRESTAPI(curriculumHandler, semesterHandler, lessonTypeHandler, disciplineHandler,
		lessonTypeAssignmentHandler, semesterDisciplineHandler, databaseHandler, []byte(jwtSecret))

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
