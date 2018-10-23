package igcsdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	igc "github.com/marni/goigc"
)

// GetRequest is giving back a json answer from requests
func GetRequest(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")
	// time when request happends

	// making a Serverinfo object that shows the uptime, info and version
	// timeFormatter is a function under igcServerUpTime
	g := ServerInfo{timeFormatter(), "Service for IGC tracks", "v1"}

	// handles the /igcinfo/
	if parts[1] == "igcinfo" && len(parts) != 1 {
		if parts[2] == "" {
			http.Error(w, "Bad request, please try again", 404)
		} else if parts[2] != "api" {
			// handles /igcinfo/<rubbish>
			http.Error(w, "something is wrong with the url path", 404)
		} else {
			// give a return a json with the service information
			// handles /igcinfo/api
			if len(parts) == 3 {
				json.NewEncoder(w).Encode(g)
			} else if parts[3] == "igc" {
				if len(parts) == 4 {
					// handles igcinfo/api/igc
					allURL(w, GlobalDb)
				} else if len(parts) == 5 {
					//handles igcinfo/api/igc/<id>
					number, err := strconv.Atoi(parts[4])
					if err != nil {
						http.Error(w, "This is not possible to convert", http.StatusBadRequest)
					}
					oneURL(w, GlobalDb, number)
				} else if len(parts) == 6 {
					//handles igcinfo/api/igc/<id>/<field>
					number, err := strconv.Atoi(parts[4])
					if err != nil {
						http.Error(w, "This is not a number", 404)
					}
					oneURLOneField(w, GlobalDb, number, parts[5])
				}
			}
		}
	} else {
		http.Error(w, "This page doesnt exsist", 404)
	}
}

//PostRequest needs to be exported for use in main and this handle all post with urls
func PostRequest(w http.ResponseWriter, r *http.Request) {
	var i IgcURL

	_, ok := GlobalDb.GetURL(i.ID)
	if ok {
		http.Error(w, "URL already exists. Use PUT to modify.", 405)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		http.Error(w, "Igc POST request must have a JSON body", http.StatusBadRequest)
		return
	}

	addToGlobalStorage(i)
	return
}

// takes all urls and put them in a slice from the map
func allURL(w http.ResponseWriter, db IgcURLStorage) {
	if db.Countl() == 0 {
		json.NewEncoder(w).Encode([]IgcURL{})
	} else {
		a := make([]IgcURL, 0, db.Countl())
		for _, s := range db.GetAll() {
			a = append(a, s)
		}
		json.NewEncoder(w).Encode(a)
	}
}

// this handles when a id is applied we get back the info about a track
func oneURL(w http.ResponseWriter, db IgcURLStorage, id int) {
	URL, ok := db.GetURL(id)
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	// else handle the url
	track, err := igc.ParseLocation(URL.URL)

	if err != nil {
		fmt.Println("Problems reading the track from URL", err)
		return
	}
	totalDistance := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}
	newTrack := igcTrack{track.Date.String(), track.Pilot, track.GliderType, track.GliderID, totalDistance}
	json.NewEncoder(w).Encode(newTrack)
}

// handles one url and the id with the field
func oneURLOneField(w http.ResponseWriter, db IgcURLStorage, id int, field string) {
	URL, ok := db.GetURL(id)
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	// else handle the url
	track, err := igc.ParseLocation(URL.URL)

	if err != nil {
		fmt.Println("Problems reading the track from URL", err)
		return
	}
	totalDistance := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}

	array := igcTrack{track.Date.String(), track.Pilot, track.GliderType, track.GliderID, totalDistance}

	// checking what the field is requested
	switch strings.ToUpper(field) {
	case "PILOT":
		json.NewEncoder(w).Encode(array.Pilot)
	case "GLIDERTYPE":
		json.NewEncoder(w).Encode(array.Glider)
	case "GLIDERID":
		json.NewEncoder(w).Encode(array.GliderID)
	case "DATE":
		json.NewEncoder(w).Encode(array.HDate)
	case "TRACKLENGTH":
		json.NewEncoder(w).Encode(array.TrackLength)

	}
}

func addToGlobalStorage(i IgcURL) {
	// this function assign
	i.ID = idForURL
	GlobalDb.Add(i)
	idForURL++

}
