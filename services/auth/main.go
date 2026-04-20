package main

import (
	"context"
	"log"
	"time"

	"github.com/IrusHunter/duckademic/services/auth/repositories"
	resthandlers "github.com/IrusHunter/duckademic/services/auth/rest_handlers"
	"github.com/IrusHunter/duckademic/services/auth/services"
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

	permissionRepository := repositories.NewPermissionRepository(database)
	roleRepository := repositories.NewRoleRepository(database)
	serviceRepository := repositories.NewServiceRepository(database)
	rolePermissionsRepository := repositories.NewRolePermissionsRepository(database)
	servicePermissionsRepository := repositories.NewServicePermissionsRepository(database)
	userRepository := repositories.NewUserRepository(database)

	adminRole := envutil.GetStringFromENV("SUPER_ADMIN_ROLE")
	if adminRole == "" {
		log.Fatalf("SUPER_ADMIN_ROLE not specified in the .env file")
	}
	roleService, adminRoleID := services.NewRoleService(roleRepository, adminRole)
	rolePermissionsService := services.NewRolePermissionsService(rolePermissionsRepository, permissionRepository,
		roleRepository)
	permissionService := services.NewPermissionService(permissionRepository, roleRepository, rolePermissionsService,
		eventBus, adminRoleID)

	serviceService := services.NewServiceService(serviceRepository)
	servicePermissionsService := services.NewServicePermissionsService(servicePermissionsRepository, serviceRepository,
		permissionRepository)

	defaultPassword := envutil.GetStringFromENV("DEFAULT_PASSWORD")
	if defaultPassword == "" {
		log.Fatalf("DEFAULT_PASSWORD not specified in the .env file")
	}
	adminLogin := envutil.GetStringFromENV("SUPER_ADMIN_LOGIN")
	if adminLogin == "" {
		log.Fatalf("SUPER_ADMIN_LOGIN not specified in the .env file")
	}
	adminPassword := envutil.GetStringFromENV("SUPER_ADMIN_PASSWORD")
	if adminPassword == "" {
		log.Fatalf("SUPER_ADMIN_PASSWORD not specified in the .env file")
	}
	jwtSecret := envutil.GetStringFromENV("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET not specified in the .env file")
	}
	userService := services.NewUserService(userRepository, roleRepository, rolePermissionsRepository, eventBus, defaultPassword,
		[]byte(jwtSecret), adminLogin, adminPassword, adminRole)

	permissionHandler := resthandlers.NewPermissionHandler(permissionService)
	roleHandler := resthandlers.NewRoleHandler(roleService)
	serviceHandler := resthandlers.NewServiceHandler(serviceService)
	rolePermissionsHandler := resthandlers.NewRolePermissionsHandler(rolePermissionsService)
	servicePermissionsHandler := resthandlers.NewServicePermissionsHandler(servicePermissionsService)
	userHandler := resthandlers.NewUserHandler(userService)
	databaseHandler := resthandlers.NewDatabaseHandler(permissionService, roleService, rolePermissionsService,
		serviceService, servicePermissionsService, userService)

	restapi := NewRESTAPI(permissionHandler, roleHandler, rolePermissionsHandler, serviceHandler,
		servicePermissionsHandler, userHandler, databaseHandler, []byte(jwtSecret))

	go func() {
		time.Sleep(1 * time.Second)
		ctx := contextutil.SetTraceID(context.Background())
		err := eventBus.PublishAccessPermissions(ctx, BuildAccessPermissions())
		if err != nil {
			log.Fatalf("Can't publish access permissions: %s", err)
		}
	}()
	err = restapi.Run(port)
	log.Fatal(err)
}
