package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

var ListenBind = "0.0.0.0"
var HttpPort = 8000
var db *sqlx.DB = nil

const skipSegmentsQuery = `SELECT "videoID",
       	"hashedVideoID" as hash,
       	"startTime",
       	"endTime",
		"UUID",
        category,
        "actionType",
       	locked,
       	votes,
       	"videoDuration",
       	"userID",
       	"description"
FROM public.sponsor_times
WHERE starts_with("hashedVideoID", $1)
  AND "votes" >= 1
  AND "category" IN
      ('sponsor', 'intro', 'outro', 'interaction', 'selfpromo', 'music_offtopic', 'preview', 'poi_highlight',
       'exclusive_access')
ORDER BY "startTime" LIMIT 100`

type SkipSegmentRow struct {
	VideoID string `db:"videoID"`
	Hash    string

	StartTime     float64 `db:"startTime"`
	EndTime       float64 `db:"endTime"`
	UUID          string  `db:"UUID"`
	Category      string
	ActionType    string `db:"actionType"`
	Locked        int    `db:"locked"`
	Votes         int    `db:"votes"`
	VideoDuration int    `db:"videoDuration"`
	UserID        string `db:"userID"`
	Description   string `db:"description"`
}

func (r SkipSegmentRow) ToSkipSegment() *SkipSegment {
	return &SkipSegment{
		Segment:       [2]float64{r.StartTime, r.EndTime},
		UUID:          r.UUID,
		Category:      r.Category,
		ActionType:    r.ActionType,
		Locked:        r.Locked,
		Votes:         r.Votes,
		VideoDuration: r.VideoDuration,
		UserID:        r.UserID,
		Description:   r.Description,
	}
}

type VideoSkipSegmentsResult struct {
	VideoID  string         `json:"videoID"`
	Hash     string         `json:"hash"`
	Segments []*SkipSegment `json:"segments"`
}

type SkipSegment struct {
	Segment       [2]float64 `json:"segment"`
	UUID          string     `json:"UUID"`
	Category      string     `json:"category"`
	ActionType    string     `json:"actionType"`
	Locked        int        `json:"locked"`
	Votes         int        `json:"votes"`
	VideoDuration int        `json:"videoDuration"`
	UserID        string     `json:"userID"`
	Description   string     `json:"description"`
}

type SkipSegmentsVideoMap map[string]*VideoSkipSegmentsResult

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
	{
		var err error
		postgresDSN := os.Getenv("POSTGRES_DSN")
		if postgresDSN == "" {
			log.Fatal("Missing or empty POSTGRES_DSN environment variable")
		}

		db, err = sqlx.Open("postgres", postgresDSN)
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
	router.HandleFunc(`/api/skipSegments/{shaPrefix:\w{4,32}}`, apiSkipSegmentsEndpoint).Methods("GET")
	addrStr := fmt.Sprintf("%s:%d", ListenBind, HttpPort)
	log.Printf("Serving requests on: '%s'", addrStr)
	log.Fatal(http.ListenAndServe(addrStr, router))
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "running")
}

func apiSkipSegmentsEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	shaPrefix := mux.Vars(r)["shaPrefix"]
	if shaPrefix == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	rows, err := db.Queryx(skipSegmentsQuery, shaPrefix)
	if err != nil {
		log.Errorf("Error getting skip segment: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	segmentsMap := SkipSegmentsVideoMap{}
	{
		parsedRowsCount := 0
		for rows.Next() {
			log.Debugf("D: %#v", rows)
			segmentRow := SkipSegmentRow{}

			err := rows.StructScan(&segmentRow)
			if err != nil {
				log.Errorf("Error while reading skipSegment rows: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if _, ok := segmentsMap[segmentRow.VideoID]; ok == false {
				// no segments for this video yet
				segmentsMap[segmentRow.VideoID] = &VideoSkipSegmentsResult{
					VideoID:  segmentRow.VideoID,
					Hash:     segmentRow.Hash,
					Segments: []*SkipSegment{segmentRow.ToSkipSegment()}[:],
				}
			} else {
				segmentsMap[segmentRow.VideoID].Segments = append(
					segmentsMap[segmentRow.VideoID].Segments, segmentRow.ToSkipSegment())
			}
			parsedRowsCount++
		}
		log.WithField("shaPrefix", shaPrefix).Debugf("skipSegments: Parsed %d rows", parsedRowsCount)
	}

	segmentTotalCount := 0
	// flatten
	results := []*VideoSkipSegmentsResult{}[:]
	for _, segResult := range segmentsMap {
		segmentTotalCount += len(segResult.Segments)
		results = append(results, segResult)
	}

	log.WithField("shaPrefix", shaPrefix).Debugf(
		"skipSegments: Returning %d videos with %d total segments", len(results), segmentTotalCount)

	jsonData, err := json.Marshal(results)
	if err != nil {
		log.Errorf("Error mashalling data into json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
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
