package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"outdoorsy/dao"
	"time"

	"github.com/gorilla/mux"
)

const (
	servicePath = "/api"
	httpPort    = ":80"
)

var daoObj *dao.DAO

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			return
		}

		// usually do auth stuff here, bearer token validation etc
		// set CORS headers if required

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*5))
		defer cancel()
		ctx = context.WithValue(ctx, "dao", daoObj)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func main() {

	log.Println("Launching demo version", os.Getenv("VERSION"))
	log.Println("Listening for http on " + httpPort)

	daoObj = dao.NewDAO()

	rtr := mux.NewRouter()
	apiPath := rtr.PathPrefix(servicePath).Subrouter()
	apiPath.Use(middleware)
	http.Handle("/", rtr)
	log.Fatal(http.ListenAndServe(httpPort, nil))
}
