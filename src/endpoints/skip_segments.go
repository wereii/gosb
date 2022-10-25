package endpoints

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gosb/src/global"
	"net/http"
)

type SkipSegmentsVideoMap map[string]*VideoSkipSegmentsResult

const SkipSegmentsQuery = `SELECT "videoID",
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
WHERE service = 'YouTube'
  AND "hashedVideoID" ILIKE $1 || '%'
  AND votes >= 0
  AND category IN ('sponsor', 'intro', 'outro', 'interaction', 'selfpromo', 
    'music_offtopic', 'preview', 'poi_highlight', 'exclusive_access')
ORDER BY "startTime" LIMIT 250`

func ApiSkipSegmentsEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	shaPrefix := mux.Vars(r)["shaPrefix"]
	if shaPrefix == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	rows, err := global.DB.Queryx(SkipSegmentsQuery, shaPrefix)
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

type SkipSegmentRow struct {
	VideoID string `db:"videoID"`
	Hash    string

	StartTime     float64 `db:"startTime"`
	EndTime       float64 `db:"endTime"`
	UUID          string  `db:"UUID"`
	Category      string
	ActionType    string  `db:"actionType"`
	Locked        int     `db:"locked"`
	Votes         int     `db:"votes"`
	VideoDuration float64 `db:"videoDuration"`
	UserID        string  `db:"userID"`
	Description   string  `db:"description"`
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
	VideoDuration float64    `json:"videoDuration"`
	UserID        string     `json:"userID"`
	Description   string     `json:"description"`
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
