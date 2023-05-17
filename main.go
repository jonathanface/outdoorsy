package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"outdoorsy/api"
	"outdoorsy/daos"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

const (
	servicePath = "/api"
	httpPort    = ":80"
)

var dao *daos.DAO

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			return
		}

		// usually do auth stuff here, bearer token validation etc
		// set CORS headers if required

		// set a timeout on api requests and attach the dao to ctx
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*5))
		defer cancel()
		ctx = context.WithValue(ctx, "dao", dao)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func main() {

	log.Println("Launching demo version", os.Getenv("VERSION"))
	log.Println("Listening for http on " + httpPort)

	dao = daos.NewDAO()
	_, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println(" received cancel signal")
		cancel()
		dao.CloseDB()
		os.Exit(1)
	}()

	rtr := mux.NewRouter()
	apiPath := rtr.PathPrefix(servicePath).Subrouter()

	// GETs
	apiPath.HandleFunc("/rentals/{rentalID}", api.RentalEndPoint).Methods("GET", "OPTIONS")

	apiPath.Use(middleware)
	http.Handle("/", rtr)
	log.Fatal(http.ListenAndServe(httpPort, nil))
}
