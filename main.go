package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"outdoorsy/api"
	"outdoorsy/daos"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const (
	servicePath = "/api"
	httpPort    = ":80"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.Info("Launching demo version: ", os.Getenv("VERSION"))
	logger.Info("Listening for HTTP on", httpPort)

	dao, err := daos.NewDAO()
	if err != nil {
		logger.Fatal("Failed to create DAO:", err)
	}

	// Graceful shutdown
	server := &http.Server{
		Addr:    httpPort,
		Handler: createHandler(dao),
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		logger.Info("Received interrupt signal. Shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			logger.Error("Error shutting down server:", err)
		}

		dao.CloseDB()

		logger.Info("Server gracefully stopped")
		os.Exit(0)
	}()

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal("Server error:", err)
	}
}

func createHandler(dao *daos.DAO) http.Handler {
	r := mux.NewRouter().PathPrefix(servicePath).Subrouter()

	// GETs
	r.HandleFunc("/rentals/{rentalID}", api.RentalEndPoint).Methods("GET", "OPTIONS")
	r.HandleFunc("/rentals", api.MultiRentalEndPoint).Methods("GET", "OPTIONS")

	r.Use(middleware(dao))

	return r
}

func middleware(dao *daos.DAO) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()

			ctx = context.WithValue(ctx, "dao", dao)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
