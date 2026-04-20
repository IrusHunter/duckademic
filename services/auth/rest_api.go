package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/auth/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/events"
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
	jwtSecret []byte,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:            platform.NewRESTAPIHelperWithAuth("RESTAPI", jwtSecret),
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
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.permissionHandler.GetAll, []string{"auth.permission"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.permissionHandler.Add, []string{"auth.permission"}),
	})
	ra.NewRoute("/permission/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.permissionHandler.Find, []string{"auth.permission"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.permissionHandler.Delete, []string{"auth.permission"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.permissionHandler.Update, []string{"auth.permission"}),
	})

	ra.NewRoute("/roles", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.roleHandler.GetAll, []string{"auth.role"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.roleHandler.Add, []string{"auth.role"}),
	})
	ra.NewRoute("/role/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.roleHandler.Find, []string{"auth.role"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.roleHandler.Delete, []string{"auth.role"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.roleHandler.Update, []string{"auth.role"}),
	})

	ra.NewRoute("/role-permissions", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.rolePermissionsHandler.GetAll, []string{"auth.role_permission"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.rolePermissionsHandler.Add, []string{"auth.role_permission"}),
	})
	ra.NewRoute("/role-permission/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.rolePermissionsHandler.Find, []string{"auth.role_permission"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.rolePermissionsHandler.Delete, []string{"auth.role_permission"}),
	})

	ra.NewRoute("/services", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.serviceHandler.GetAll, []string{"auth.service"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.serviceHandler.Add, []string{"auth.service"}),
	})
	ra.NewRoute("/service/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.serviceHandler.Find, []string{"auth.service"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.serviceHandler.Delete, []string{"auth.service"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.serviceHandler.Update, []string{"auth.service"}),
	})

	ra.NewRoute("/service-permissions", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.servicePermissionHandler.GetAll, []string{"auth.service_permission"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.servicePermissionHandler.Add, []string{"auth.service_permission"}),
	})
	ra.NewRoute("/service-permission/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.servicePermissionHandler.Find, []string{"auth.service_permission"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.servicePermissionHandler.Delete, []string{"auth.service_permission"}),
	})

	ra.NewRoute("/users", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.userHandler.GetAll, []string{"auth.user"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.userHandler.Add, []string{"auth.user"}),
	})
	ra.NewRoute("/user/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.userHandler.Find, []string{"auth.user"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.userHandler.Update, []string{"auth.user"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.userHandler.Delete, []string{"auth.user"}),
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.userHandler.Login)(r.Context(), w, r)
	})
	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandler(ra.userHandler.Refresh)(r.Context(), w, r)
	})
	http.HandleFunc("/reset-password/{id}", func(w http.ResponseWriter, r *http.Request) {
		ra.NewDefaultHandlerWithAuth(ra.userHandler.ResetPassword, []string{"auth.user.reset_password"})(r.Context(), w, r)
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
func BuildAccessPermissions() []events.AccessPermissionRE {
	return []events.AccessPermissionRE{
		// permissions
		{Name: "auth.permission"},

		// roles
		{Name: "auth.role"},

		// role-permissions
		{Name: "auth.role_permission"},

		// services
		{Name: "auth.service"},

		// service-permissions
		{Name: "auth.service_permission"},

		// users
		{Name: "auth.user"},
		{Name: "auth.user.reset_password"},
	}
}
