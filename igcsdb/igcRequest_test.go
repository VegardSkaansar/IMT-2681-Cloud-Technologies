package igcsdb

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marni/goigc"
)

func Test_getRequest_malformedURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(GetRequest))
	defer ts.Close()

	// malformed Url
	resp, err := http.Get(ts.URL + "/igcinfo/ape")
	if err != nil {
		t.Error("Error with get request", err)
	}
	if resp.StatusCode != 404 {
		t.Errorf("For get %s, expected StatusCode %d, received %d", "/igcinfo/ape",
			404, resp.StatusCode)
		return
	}

}

func Test_getRequest_serviceInfo(t *testing.T) {
	ServerStart()

	ts := httptest.NewServer(http.HandlerFunc(GetRequest))
	defer ts.Close()

	testService := ServerInfo{timeFormatter(), "Service for IGC tracks", "v1"}

	resp, err := http.Get(ts.URL + "/igcinfo/api")

	if err != nil {
		t.Error("Error making the get request")
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}

	json.NewDecoder(resp.Body).Decode(&testService)
}

func Test_getRequest_AllURL(t *testing.T) {

	GlobalDb = &IgcURLDB{}
	GlobalDb.Init()
	testIgc := IgcURL{"http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc", 0}
	testIgc2 := IgcURL{"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc", 1}
	GlobalDb.Add(testIgc)
	GlobalDb.Add(testIgc2)

	ts := httptest.NewServer(http.HandlerFunc(GetRequest))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/igcinfo/api/igc")

	if err != nil {
		t.Error("Error making the get request")
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}

	var a []IgcURL
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got error: %s", err)
	}

	if len(a) != 2 {
		t.Errorf("Excpected array with one element, got %v", a)
	}
	// could have made this to a array with this vaule for to add more igc later but here
	// wont change anything in the future, thtas why im hardcoding
	if a[0].URL != testIgc.URL || a[0].ID != testIgc.ID {
		t.Errorf("URL IS not the same. Got: %v, Expected: %v\n", a[0], testIgc)
	}

	if a[1].URL != testIgc2.URL || a[1].ID != testIgc2.ID {
		t.Errorf("URL IS not the same. Got: %v, Expected: %v\n", a[0], testIgc)
	}

}

func Test_getRequest_OneURL(t *testing.T) {

	GlobalDb = &IgcURLDB{}
	GlobalDb.Init()
	testIgc := IgcURL{"http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc", 0}
	GlobalDb.Add(testIgc)

	ts := httptest.NewServer(http.HandlerFunc(GetRequest))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/igcinfo/api/igc/0")

	if err != nil {
		t.Error("Error making the get request")
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}
	var i igcTrack
	err = json.NewDecoder(resp.Body).Decode(&i)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got error: %s", err)
	}
	// hardcoding the function from request to check if it gives the right values
	track, _ := igc.ParseLocation(testIgc.URL)
	totalDistance := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}
	newTrack := igcTrack{track.Date.String(), track.Pilot, track.GliderType, track.GliderID, totalDistance}

	if i.Pilot != newTrack.Pilot || i.Glider != newTrack.Glider || i.GliderID != newTrack.GliderID || i.HDate != newTrack.HDate || i.TrackLength != newTrack.TrackLength {
		t.Error("The tracker JSON doesnt match with the URL acctually track")
	}
}
