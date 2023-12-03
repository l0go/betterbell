package main

import (
	"codeberg.org/logo/betterbell/internal"

	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/go-co-op/gocron"
)

var tpl *template.Template
var jobs internal.JobsState

func main() {
	log.Println("Running")

	// Create database
	db, err := sql.Open("sqlite", "./bell.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create jobs table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS CronJobs (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, expression TEXT, enabled BOOLEAN);")
	if err != nil {
		log.Fatal(err)
	}

	// Create user table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, hash TEXT, salt TEXT);")
	if err != nil {
		log.Fatal(err)
	}

	// Create the gocron scheduler
	scheduler := gocron.NewScheduler(time.Local)
	jobs = internal.JobsState{
		Scheduler: scheduler,
		DB:        db,
	}
	jobs.Persist()

	// HTTP
	tpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/login", getLogin)
	http.HandleFunc("/job/add", getJobAdd)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err = http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Index reachable at /
// Contains a list of the jobs and actions that can be acted upon them
func getRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		action, id, _ := strings.Cut(r.FormValue("action"), "-")
		switch action {
		case "ring":
			if err := internal.Ring(); err != nil {
				log.Fatal(err)
			}
		case "delete":
			if err := jobs.Remove(id); err != nil {
                log.Printf("Removing job failed: %s\n", err)
			}
        case "toggle":
            if err := jobs.Toggle(id); err != nil {
                log.Printf("Toggling job failed: %s\n", err)
            }
        }
	}
	tpl.ExecuteTemplate(w, "index.html", jobs.Get())
}

// Login page reachable at /login
// Allows you to login to access the rest of the service
func getLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "templates/login.html")
	}
}

// Job adding page reachable at /job/add
// Allows you to add a job
func getJobAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		expression := r.FormValue("cron-expression")

		if len(title) == 0 || len(expression) == 0 {
			log.Fatal("Job form data is invalid")
		}

		if err := jobs.Add(title, expression); err != nil {
			log.Fatal(err)
		}
	}
	tpl.ExecuteTemplate(w, "add-job.html", nil)
}
