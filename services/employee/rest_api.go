package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	resthandlers "github.com/IrusHunter/duckademic/services/employees/rest_handlers"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
)

// RESTAPI represents a RESTful HTTP server that can be started on a given port.
type RESTAPI interface {
	Run(int) error // Run starts the REST API server on the specified port.
}

// NewRESTAPI creates a new RESTAPI instance.
//
// It requires the academic rank (arh) and database (dh) handlers.
func NewRESTAPI(arh resthandlers.AcademicRankHandler, dh resthandlers.DatabaseHandler) RESTAPI {
	return &restapi{academicRankHandler: arh, databaseHandler: dh}
}

type restapi struct {
	academicRankHandler resthandlers.AcademicRankHandler
	databaseHandler     resthandlers.DatabaseHandler
}

func (ra *restapi) Run(port int) error {
	newHandler("/academic-ranks", ra.academicRanksRouter)
	newHandler("/academic-rank/{id}", ra.academicRankIDRouter)
	newHandler("/seed", ra.databaseHandler.Seed)

	log.Printf("Server start at port %d \n", port)

	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func (ra *restapi) academicRankIDRouter(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ra.academicRankHandler.Find(ctx, w, r)
		return
	case http.MethodDelete:
		ra.academicRankHandler.Delete(ctx, w, r)
		return
	case http.MethodPut:
		ra.academicRankHandler.Update(ctx, w, r)
		return
	}

	jsonutil.ResponseWithError(
		w, 405, fmt.Errorf("method %q not available (available methods GET, POST)", r.Method),
	)
}

func (ra *restapi) academicRanksRouter(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ra.academicRankHandler.GetAll(ctx, w, r)
		return
	case http.MethodPost:
		ra.academicRankHandler.Add(ctx, w, r)
		return
	}

	jsonutil.ResponseWithError(
		w, 405, fmt.Errorf("method %q not available (available methods GET, POST)", r.Method),
	)
}

func newHandler(path string, f func(context.Context, http.ResponseWriter, *http.Request),
) {
	middleF := func(w http.ResponseWriter, r *http.Request) {
		ctx := contextutil.SetTraceID(context.Background())
		f(ctx, w, r)
	}
	http.HandleFunc(path, middleF)
}
