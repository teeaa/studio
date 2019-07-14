package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/teeaa/studio/internal/bookings"
	"github.com/teeaa/studio/internal/classes"
)

func startServer() *http.Server {
	log.Info("Starting REST API")
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      getRouter(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("Unable to start REST API: ", err)
		}
	}()

	return srv
}

func getRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(false)
	router.Use(logRequest, setHeaders)

	gormDB := Connect()
	classes.Routes(gormDB, router.PathPrefix("/classes").Subrouter())
	bookings.Routes(gormDB, router.PathPrefix("/bookings").Subrouter())

	return router
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("%s: %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func setHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
