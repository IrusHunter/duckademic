package main

import (
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/asset/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/platform"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

func NewRESTAPI(
	ch resthandlers.ClassroomHandler,
	dh resthandlers.DatabaseHandler,
	jwtSecrete []byte,
) RESTAPI {
	return &restapi{
		RESTAPIHelper:    platform.NewRESTAPIHelperWithAuth("RESTAPI", jwtSecrete),
		classroomHandler: ch,
		databaseHandler:  dh,
	}
}

type restapi struct {
	platform.RESTAPIHelper
	classroomHandler resthandlers.ClassroomHandler
	databaseHandler  resthandlers.DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	ra.NewRoute("/classrooms", map[string]platform.HandlerFunc{
		http.MethodGet:  ra.NewDefaultHandlerWithAuth(ra.classroomHandler.GetAll, []string{"asset.classroom"}),
		http.MethodPost: ra.NewDefaultHandlerWithAuth(ra.classroomHandler.Add, []string{"asset.classroom"}),
	})
	ra.NewRoute("/classroom/{id}", map[string]platform.HandlerFunc{
		http.MethodGet:    ra.NewDefaultHandlerWithAuth(ra.classroomHandler.Find, []string{"asset.classroom"}),
		http.MethodDelete: ra.NewDefaultHandlerWithAuth(ra.classroomHandler.Delete, []string{"asset.classroom"}),
		http.MethodPut:    ra.NewDefaultHandlerWithAuth(ra.classroomHandler.Update, []string{"asset.classroom"}),
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
		{Name: "auth.classroom"},
	}
}
