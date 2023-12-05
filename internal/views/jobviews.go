package views

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"codeberg.org/logo/betterbell/internal"
	"github.com/go-co-op/gocron"
)

type JobViews struct {
	State          internal.JobsState
	Peers          *PeerViews
	IndexTemplate  *template.Template
	AddJobTemplate *template.Template
}

func CreateJobs(db *sql.DB, peers *PeerViews) JobViews {
	var v JobViews

	scheduler := gocron.NewScheduler(time.Local)
	v.State = internal.JobsState{
		Scheduler: scheduler,
		Peers:     &peers.State,
		DB:        db,
	}
	v.State.Persist()

	v.Peers = peers

	v.IndexTemplate = template.Must(template.ParseFiles("templates/index.html", "templates/base.html"))
	v.AddJobTemplate = template.Must(template.ParseFiles("templates/add-job.html", "templates/base.html"))

	return v
}

// Index reachable at /
// Contains a list of the jobs and actions that can be acted upon them
func (v JobViews) GetIndex(w http.ResponseWriter, r *http.Request) {
	// Check if we have permission
	v.Peers.Auth.CheckPermission(w, r)

	if r.Method == http.MethodPost {
		action, id, _ := strings.Cut(r.FormValue("action"), "-")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("Could not convert id to int: %s\n", err)
		}

		switch action {
		case "ring":
			if err := internal.Ring(idInt, &v.Peers.State); err != nil {
				log.Fatal(err)
			}
		case "delete":
			if err := v.State.Remove(idInt); err != nil {
				log.Printf("Removing job failed: %s\n", err)
			}
		case "toggle":
			if err := v.State.Toggle(idInt); err != nil {
				log.Printf("Toggling job failed: %s\n", err)
			}
		}
	}
	v.IndexTemplate.ExecuteTemplate(w, "base", v.State.Get())
}

// Job adding page reachable at /job/add
// Allows you to add a job
func (v JobViews) GetJobAdd(w http.ResponseWriter, r *http.Request) {
	// Check if we have permission
	v.Peers.Auth.CheckPermission(w, r)

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		expression := r.FormValue("cron-expression")

		if len(title) == 0 || len(expression) == 0 {
			log.Fatal("Job form data is invalid")
		}

		if err := v.State.Add(title, expression); err != nil {
			log.Fatal(err)
		}
	}
	v.AddJobTemplate.ExecuteTemplate(w, "base", nil)
}
