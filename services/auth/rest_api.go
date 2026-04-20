package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/auth/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

func NewRESTAPI(
	ph resthandlers.PermissionHandler,
	rh resthandlers.RoleHandler,
	rph resthandlers.RolePermissionsHandler,
	sh resthandlers.ServiceHandler,
	sph resthandlers.ServicePermissionsHandler,
	uh resthandlers.UserHandler,
	dh resthandlers.DatabaseHandler,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:            platform.NewRESTAPIHelper("RESTAPI"),
		permissionHandler:        ph,
		roleHandler:              rh,
		rolePermissionsHandler:   rph,
		serviceHandler:           sh,
		servicePermissionHandler: sph,
		userHandler:              uh,
		databaseHandler:          dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	permissionHandler        resthandlers.PermissionHandler
	roleHandler              resthandlers.RoleHandler
	rolePermissionsHandler   resthandlers.RolePermissionsHandler
	serviceHandler           resthandlers.ServiceHandler
	servicePermissionHandler resthandlers.ServicePermissionsHandler
	userHandler              resthandlers.UserHandler
	databaseHandler          resthandlers.DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/permissions", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.permissionHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.permissionHandler.Add),
	})
	ra.NewRoute("/permission/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.permissionHandler.Find),
		http.MethodDelete: ra.NewDefaultHandler(ra.permissionHandler.Delete),
		http.MethodPut:    ra.NewDefaultHandler(ra.permissionHandler.Update),
	})

	ra.NewRoute("/roles", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.roleHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.roleHandler.Add),
	})
	ra.NewRoute("/role/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.roleHandler.Find),
		http.MethodDelete: ra.NewDefaultHandler(ra.roleHandler.Delete),
		http.MethodPut:    ra.NewDefaultHandler(ra.roleHandler.Update),
	})

	ra.NewRoute("/role-permissions", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.rolePermissionsHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.rolePermissionsHandler.Add),
	})
	ra.NewRoute("/role-permission/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.rolePermissionsHandler.Find),
		http.MethodDelete: ra.NewDefaultHandler(ra.rolePermissionsHandler.Delete),
	})

	ra.NewRoute("/services", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.serviceHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.serviceHandler.Add),
	})
	ra.NewRoute("/service/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.serviceHandler.Find),
		http.MethodDelete: ra.NewDefaultHandler(ra.serviceHandler.Delete),
		http.MethodPut:    ra.NewDefaultHandler(ra.serviceHandler.Update),
	})

	ra.NewRoute("/service-permissions", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandler(ra.servicePermissionHandler.GetAll),
		http.MethodPost: ra.NewDefaultHandler(ra.servicePermissionHandler.Add),
	})
	ra.NewRoute("/service-permission/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandler(ra.servicePermissionHandler.Find),
		http.MethodDelete: ra.NewDefaultHandler(ra.servicePermissionHandler.Delete),
	})

	ra.NewRoute("/users", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.userHandler.GetAll),
	})
	ra.NewRoute("/user/{id}", map[string]platform.HandlerFunc{
		http.MethodGet: ra.NewDefaultHandler(ra.userHandler.Find),
		http.MethodPut: ra.NewDefaultHandler(ra.userHandler.Update),
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.userHandler.Login)(r.Context(), w, r)
	})
	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.userHandler.Refresh)(r.Context(), w, r)
	})
	http.HandleFunc("/reset-password/{id}", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.userHandler.ResetPassword)(r.Context(), w, r)
	})
	http.HandleFunc("/change-password", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.userHandler.ChangePassword)(r.Context(), w, r)
	})

	http.HandleFunc("/seed", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.databaseHandler.Seed)(r.Context(), w, r)
	})
	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.databaseHandler.Clear)(r.Context(), w, r)
	})

	log.Printf("Server start at port %d \n", port)

	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
