package main

import (
	"context"
	"log"
	"time"

	"github.com/IrusHunter/duckademic/services/course/repositories"
	resthandlers "github.com/IrusHunter/duckademic/services/course/rest_handlers"
	"github.com/IrusHunter/duckademic/services/course/services"
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
	teacherRepository := repositories.NewTeacherRepository(database)
	courseRepository := repositories.NewCourseRepository(database)
	studentCourseRepository := repositories.NewStudentCourseRepository(database)
	teacherCourseRepository := repositories.NewTeacherCourseRepository(database)
	taskRepository := repositories.NewTaskRepository(database)
	taskStudentRepository := repositories.NewTaskStudentRepository(database)

	studentService := services.NewStudentService(studentRepository, eventBus)
	teacherService := services.NewTeacherService(teacherRepository, eventBus)
	courseService := services.NewCourseService(courseRepository, teacherRepository, eventBus)
	studentCourseService := services.NewStudentCourseService(studentCourseRepository, studentRepository, courseRepository)
	teacherCourseService := services.NewTeacherCourseService(teacherCourseRepository, teacherRepository, courseRepository)
	taskService := services.NewTaskService(taskRepository, courseRepository)
	taskStudentService := services.NewTaskStudentService(taskStudentRepository, taskRepository, studentRepository)

	studentHandler := resthandlers.NewStudentHandler(studentService)
	teacherHandler := resthandlers.NewTeacherHandler(teacherService)
	courseHandler := resthandlers.NewCourseHandler(courseService)
	studentCourseHandler := resthandlers.NewStudentCourseHandler(studentCourseService)
	teacherCourseHandler := resthandlers.NewTeacherCourseHandler(teacherCourseService)
	taskHandler := resthandlers.NewTaskHandler(taskService)
	taskStudentHandler := resthandlers.NewTaskStudentHandler(taskStudentService)
	databaseHandler := resthandlers.NewDatabaseHandler(studentService, teacherService, courseService, studentCourseService,
		teacherCourseService, taskService, taskStudentService)

	jwtSecret := envutil.GetStringFromENV("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET not specified in the .env file")
	}
	restapi := NewRESTAPI(studentHandler, teacherHandler, courseHandler, studentCourseHandler, teacherCourseHandler,
		taskHandler, taskStudentHandler, databaseHandler, []byte(jwtSecret))

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
