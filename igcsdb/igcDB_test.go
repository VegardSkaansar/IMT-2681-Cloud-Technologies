package igcsdb

import (
	"fmt"
	"testing"
)

func Test_addURL(t *testing.T) {

	rangeOfID := 0

	db := &IgcURLDB{}
	igcData := IgcURL{"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc", 0}

	db.Init()
	rangeOfID++
	igcData.ID = rangeOfID
	db.Add(igcData)

	if db.Countl() != 1 {
		t.Error("This is not the right amount of URLS")
	}

	// test if the url accutally has been added form the add function
	if igcData.ID > db.Countl() {
		t.Error("a URL cant be added right since the id is out of range")
	}
	if igcData.URL != "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc" {
		t.Error("Igc file was not added.")
	}

}

func Test_moreUrls(t *testing.T) {

	testData := map[int]IgcURL{
		1: {"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc", 1},
		2: {"http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc", 2},
	}

	db := &IgcURLDB{}
	db.Init()

	for _, s := range testData {
		db.Add(s)
	}

	fmt.Println(len(db.igcUrls))
	if db.Countl() != len(testData) {
		t.Error("Wrong number of URLS")
	}
	for igcNum := range db.igcUrls {
		igc, _ := db.GetURL(igcNum)
		igcTest, _ := testData[igcNum]

		if igc.URL != igcTest.URL {
			t.Error("Wrong Url ")
		}

		if igc.ID != igcTest.ID {
			t.Error("wrong Id")
		}
	}

}
