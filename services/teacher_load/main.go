package main

import (
	"log"

	"github.com/IrusHunter/duckademic/services/teacher_load/repositories"
	resthandlers "github.com/IrusHunter/duckademic/services/teacher_load/rest_handlers"
	"github.com/IrusHunter/duckademic/services/teacher_load/services"
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
	groupCohortRepository := repositories.NewGroupCohortRepository(database)
	lessonTypeRepository := repositories.NewLessonTypeRepository(database)
	disciplineRepository := repositories.NewDisciplineRepository(database)

	teacherService := services.NewTeacherService(teacherRepository, eventBus)
	groupCohortService := services.NewGroupCohortService(groupCohortRepository, eventBus)
	lessonTypeService := services.NewLessonTypeService(lessonTypeRepository, eventBus)
	disciplineService := services.NewDisciplineService(disciplineRepository, eventBus)

	teacherHandler := resthandlers.NewTeacherHandler(teacherService)
	groupCohortHandler := resthandlers.NewGroupCohortHandler(groupCohortService)
	lessonTypeHandler := resthandlers.NewLessonTypeHandler(lessonTypeService)
	disciplineHandler := resthandlers.NewDisciplineHandler(disciplineService)
	databaseHandler := resthandlers.NewDatabaseHandler()

	restapi := NewRESTAPI(teacherHandler, groupCohortHandler, lessonTypeHandler, disciplineHandler, databaseHandler)

	err = restapi.Run(port)
	log.Fatal(err)
}
