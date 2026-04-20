package main

import (
	"context"
	"log"
	"time"

	"github.com/IrusHunter/duckademic/services/student_group/repositories"
	resthandlers "github.com/IrusHunter/duckademic/services/student_group/rest_handlers"
	"github.com/IrusHunter/duckademic/services/student_group/services"
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

	semesterRepository := repositories.NewSemesterRepository(database)
	studentRepository := repositories.NewStudentRepository(database)
	groupCohortRepository := repositories.NewGroupCohortRepository(database)
	studentGroupRepository := repositories.NewStudentGroupRepository(database)
	groupMemberRepository := repositories.NewGroupMemberRepository(database)
	disciplineRepository := repositories.NewDisciplineRepository(database)
	lessonTypeRepository := repositories.NewLessonTypeRepository(database)
	groupCohortAssignmentRepository := repositories.NewGroupCohortAssignmentRepository(database)

	semesterService := services.NewSemesterService(semesterRepository, eventBus)
	studentService := services.NewStudentService(studentRepository, eventBus)
	groupCohortService := services.NewGroupCohortService(groupCohortRepository, semesterRepository, eventBus)
	studentGroupService := services.NewStudentGroupService(studentGroupRepository, groupCohortRepository, eventBus)
	groupMemberService := services.NewGroupMemberService(groupMemberRepository, studentRepository, groupCohortRepository,
		studentGroupRepository, eventBus)
	disciplineService := services.NewDisciplineService(disciplineRepository, eventBus)
	lessonTypeService := services.NewLessonTypeService(lessonTypeRepository, eventBus)
	groupCohortAssignmentService := services.NewGroupCohortAssignmentService(groupCohortAssignmentRepository,
		groupCohortRepository, disciplineRepository, lessonTypeRepository, eventBus)

	semesterHandler := resthandlers.NewSemesterHandler(semesterService)
	studentHandler := resthandlers.NewStudentHandler(studentService)
	groupCohortHandler := resthandlers.NewGroupCohortHandler(groupCohortService)
	studentGroupHandler := resthandlers.NewStudentGroupHandler(studentGroupService)
	groupMemberHandler := resthandlers.NewGroupMemberHandler(groupMemberService)
	disciplineHandler := resthandlers.NewDisciplineHandler(disciplineService)
	lessonTypeHandler := resthandlers.NewLessonTypeHandler(lessonTypeService)
	groupCohortAssignmentHandler := resthandlers.NewGroupCohortAssignmentHandler(groupCohortAssignmentService)
	databaseHandler := resthandlers.NewDatabaseHandler(studentService, semesterService, groupCohortService,
		studentGroupService, groupMemberService, disciplineService, lessonTypeService, groupCohortAssignmentService)

	jwtSecret := envutil.GetStringFromENV("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET not specified in the .env file")
	}
	restapi := NewRESTAPI(studentHandler, semesterHandler, groupCohortHandler, studentGroupHandler, groupMemberHandler,
		lessonTypeHandler, disciplineHandler, groupCohortAssignmentHandler, databaseHandler, []byte(jwtSecret))

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
