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
		log.Fatalf("can't get port value: %s", err.Error())
	}

	database, err := db.NewDefaultConnection()
	if err != nil {
		log.Fatalf("Can't connect to database: %v", err)
	}

	upstreamRepository := NewUpstreamRepository(database)
	endpointRepository := NewEndpointRepository(upstreamRepository)

	NewUpstreamService(upstreamRepository)
	endpointService := NewEndpointService(endpointRepository)

	proxyHandler := NewProxyHandler(endpointService, http.DefaultClient)

	restapi := NewRESTAPI(proxyHandler)

	err = restapi.Run(port)
	log.Fatal(err)
}
