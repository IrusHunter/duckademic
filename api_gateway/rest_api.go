package main

import (
	"log"
	"net/http"
	"strconv"
)

type RESTAPI interface {
	Run(int) error
}

func NewRESTAPI(pc ProxyHandler) RESTAPI {
	return &restapi{proxy: pc}
}

type restapi struct {
	proxy ProxyHandler
}

func (ra *restapi) Run(port int) error {
	newHandler("/", ra.proxy.HandlePath)

	log.Printf("Server start at port %d \n", port)

	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func newHandler(path string, f http.HandlerFunc) {
	http.HandleFunc(path, corsMiddleware(f))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
