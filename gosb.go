package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

var ListenBind = "0.0.0.0"
var HttpPort = 8000

const skipSegmentsQuery = "select * from "

func getEnvOpts() {
	envHttpPort, ok := os.LookupEnv("HTTP_PORT")
	if ok {
		i, err := strconv.Atoi(envHttpPort)
		if err != nil {
			log.Fatalf("Failed converting HTTP_PORT value '%s' to number: %s", envHttpPort, err)
		}
		HttpPort = i
	}
	_, ok = os.LookupEnv("DEBUG")
	if ok {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	var db *sql.DB
	{
		var err error
		postgresDSN := os.Getenv("POSTGRES_DSN")
		if postgresDSN == "" {
			log.Fatal("Missing or empty POSTGRES_DSN environment variable")
		}

		db, err = sql.Open("postgres", postgresDSN)
		if err != nil {
			log.Fatalf("Failed seeting up db connection: %s", err)
		}

		err = db.Ping()
		if err != nil {
			log.Fatalf("Failed connecting to database: %s", err)
		}

		defer func() {
			err := db.Close()
			if err != nil {
				log.Printf("Error while closing database: %s", err)
			}
		}()
	}
	getEnvOpts()

	handleRequests()
}

func handleRequests() {
	router := mux.NewRouter()
	//router.StrictSlash()
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("Debug level enabled, logging web requests enabled")
		router.Use(loggingMiddleware)
	}
	router.HandleFunc("/", indexPage)
	router.HandleFunc(`/api/skipSegments/{shaPrefix:\w{4,32}}`, apiSkipSegmentsPage).Methods("GET")
	addrStr := fmt.Sprintf("%s:%d", ListenBind, HttpPort)
	log.Printf("Serving requests on: '%s'", addrStr)
	log.Fatal(http.ListenAndServe(addrStr, router))
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "running")
}

func apiSkipSegmentsPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "[]")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logEntry := log.WithField("remote", r.RemoteAddr)
		if r.Header.Get("X-Forwarded-For") != "" {
			logEntry.WithField("X-Forwarded-For", r.Header.Get("X-Forwarded-For"))
		}
		if r.Header.Get("X-Real-IP") != "" {
			logEntry.WithField("X-Real-IP", r.Header.Get("X-Real-IP"))
		}
		logEntry.Debugf("Request: %s", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
