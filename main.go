package main

import (
	"fmt"
	"os"

	"codeberg.org/logo/betterbell/internal"
	"codeberg.org/logo/betterbell/internal/views"

	"database/sql"
	"log"
	"net/http"

	_ "github.com/glebarez/go-sqlite"
)

func main() {
	log.Println("Running")

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3333"
	}

	dbLocation := os.Getenv("DB_LOCATION")
	if dbLocation == "" {
		dbLocation = "./bell.db"
	}

	// Connect to the database
	db, err := sql.Open("sqlite", dbLocation)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Now we can create all of the tables
	internal.CreateTables(db)

	// Initialize Views
	auth := views.CreateAuth(db)
	peers := views.CreatePeers(db, &auth)
	jobs := views.CreateJobs(db, &peers)
	ring := views.RingViews{Peers: peers}

	// HTTP
	// Define templates
	http.HandleFunc("/", jobs.GetIndex)
	http.HandleFunc("/job/add/", jobs.GetJobAdd)
	http.HandleFunc("/login/", auth.GetLogin)
	http.HandleFunc("/register/", auth.GetRegister)
	http.HandleFunc("/peers/", peers.GetPeers)
	http.HandleFunc("/peers/remove/", peers.GetPeersRemove)
	http.HandleFunc("/ring/", ring.GetRing)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err = http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil)
	if err != nil {
		log.Fatal(err)
	}
}
