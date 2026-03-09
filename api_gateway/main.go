package main

import (
	"log"
	"net/http"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/IrusHunter/duckademic/shared/envutil"
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
		log.Fatalf("Can't migrate the database^ %s", err.Error())
	}

	upstreamRepository := NewUpstreamRepository(database)
	endpointRepository := NewEndpointRepository(database, upstreamRepository)

	upstreamService := NewUpstreamService(upstreamRepository)
	endpointService := NewEndpointService(endpointRepository, upstreamRepository)

	proxyHandler := NewProxyHandler(endpointService, http.DefaultClient)
	databaseHandler := NewDatabaseHandler(upstreamService, endpointService)

	restapi := NewRESTAPI(proxyHandler, databaseHandler)

	err = restapi.Run(port)
	log.Fatal(err)
}
