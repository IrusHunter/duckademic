package main

import (
	"log"

	"github.com/IrusHunter/duckademic/services/auth/repositories"
	resthandlers "github.com/IrusHunter/duckademic/services/auth/rest_handlers"
	"github.com/IrusHunter/duckademic/services/auth/services"
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

	permissionService := services.NewPermissionService(permissionRepository, eventBus)
	roleService := services.NewRoleService(roleRepository)
	serviceService := services.NewServiceService(serviceRepository)
	rolePermissionsService := services.NewRolePermissionsService(rolePermissionsRepository, permissionRepository,
		roleRepository)
	servicePermissionsService := services.NewServicePermissionsService(servicePermissionsRepository, serviceRepository,
		permissionRepository)

	permissionHandler := resthandlers.NewPermissionHandler(permissionService)
	roleHandler := resthandlers.NewRoleHandler(roleService)
	serviceHandler := resthandlers.NewServiceHandler(serviceService)
	rolePermissionsHandler := resthandlers.NewRolePermissionsHandler(rolePermissionsService)
	servicePermissionsHandler := resthandlers.NewServicePermissionsHandler(servicePermissionsService)
	databaseHandler := resthandlers.NewDatabaseHandler(roleService, rolePermissionsService, serviceService,
		servicePermissionsService)

	restapi := NewRESTAPI(permissionHandler, roleHandler, rolePermissionsHandler, serviceHandler,
		servicePermissionsHandler, databaseHandler)

	err = restapi.Run(port)
	log.Fatal(err)
}
