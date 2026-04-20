package main

import (
	"context"
	"log"
	"time"

	"github.com/IrusHunter/duckademic/services/teacher_load/repositories"
	resthandlers "github.com/IrusHunter/duckademic/services/teacher_load/rest_handlers"
	"github.com/IrusHunter/duckademic/services/teacher_load/services"
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

	teacherRepository := repositories.NewTeacherRepository(database)
	lessonTypeRepository := repositories.NewLessonTypeRepository(database)
	disciplineRepository := repositories.NewDisciplineRepository(database)
	teacherLoadRepository := repositories.NewTeacherLoadRepository(database)

	teacherService := services.NewTeacherService(teacherRepository, eventBus)
	lessonTypeService := services.NewLessonTypeService(lessonTypeRepository, eventBus)
	disciplineService := services.NewDisciplineService(disciplineRepository, eventBus)
	teacherLoadService := services.NewTeacherLoadService(teacherLoadRepository, teacherRepository, disciplineRepository,
		lessonTypeRepository, eventBus)

	teacherHandler := resthandlers.NewTeacherHandler(teacherService)
	lessonTypeHandler := resthandlers.NewLessonTypeHandler(lessonTypeService)
	disciplineHandler := resthandlers.NewDisciplineHandler(disciplineService)
	teacherLoadHandler := resthandlers.NewTeacherLoadHandler(teacherLoadService)
	databaseHandler := resthandlers.NewDatabaseHandler(teacherLoadService, teacherService, disciplineService,
		lessonTypeService)

	jwtSecret := envutil.GetStringFromENV("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET not specified in the .env file")
	}
	restapi := NewRESTAPI(teacherHandler, lessonTypeHandler, disciplineHandler,
		teacherLoadHandler, databaseHandler, []byte(jwtSecret))

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
