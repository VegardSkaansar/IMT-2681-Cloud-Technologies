package main

import (
	"fmt"
	"goprojects/igcTracker/igcsdb"
	"log"
	"net/http"
	"os"
)

//global clock to controll the uptime

func determineListenAddress() (string, error) {
	igcsdb.ServerStart()
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func handlerIGC(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		igcsdb.GetRequest(w, r)
	case "POST":
		igcsdb.PostRequest(w, r)
	}
}

func main() {

	// Using in-memory storage
	igcsdb.GlobalDb = &igcsdb.IgcURLDB{}

	// initialsing the database
	igcsdb.GlobalDb.Init()
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/igcinfo/", handlerIGC)
	http.ListenAndServe(addr, nil)

}
