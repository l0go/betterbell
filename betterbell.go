package main

import (
	"strconv"

	"codeberg.org/logo/betterbell/internal"

	"database/sql"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/go-co-op/gocron"
)

var tpl map[string]*template.Template
var auth internal.AuthState
var jobs internal.JobsState
var peers map[string]string

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

	// Create the authentication state
	auth = internal.AuthState{
		DB:       db,
		Sessions: make(map[[16]byte]struct{}),
	}

	// HTTP
	// Define templates
	tpl = make(map[string]*template.Template)
	tpl["index"] = template.Must(template.ParseFiles("templates/index.html", "templates/base.html"))
	tpl["login"] = template.Must(template.ParseFiles("templates/login.html", "templates/base.html"))
	tpl["register"] = template.Must(template.ParseFiles("templates/register.html", "templates/base.html"))
	tpl["add-job"] = template.Must(template.ParseFiles("templates/add-job.html", "templates/base.html"))
	tpl["peers"] = template.Must(template.ParseFiles("templates/peers.html", "templates/base.html"))

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/login/", getLogin)
	http.HandleFunc("/register/", getRegister)
	http.HandleFunc("/job/add/", getJobAdd)
	http.HandleFunc("/peers/", getPeers)
	http.HandleFunc("/ring/", getRing)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err = http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Index reachable at /
// Contains a list of the jobs and actions that can be acted upon them
func getRoot(w http.ResponseWriter, r *http.Request) {
	has, err := auth.HasPermission(r)
	if err != nil {
		log.Println(err)
	}

	if !has {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		action, id, _ := strings.Cut(r.FormValue("action"), "-")
		switch action {
		case "ring":
			intId, _ := strconv.Atoi(id)
			if err := internal.Ring(intId); err != nil {
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
	tpl["index"].ExecuteTemplate(w, "base", nil)
}

// Login page reachable at /login
// Allows you to login to access the rest of the service
func getLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if len(username) == 0 || len(password) == 0 {
			log.Fatal("Job form data is invalid")
		}

		if status, err := auth.Check(username, password); status == internal.LoginSuccess && err == nil {
			// Account exists and the password is valid
			session, err := auth.GrantSession(r)
			if err != nil {
				log.Println(err)
				return
			}
			http.SetCookie(w, &session)
			log.Println(internal.FormatLoginStatus(status))
		} else {
			// Account is not valid
			log.Println("Account does not exist...")
			log.Println(internal.FormatLoginStatus(status))
		}
		return
	}
	tpl["login"].ExecuteTemplate(w, "base", nil)
}

// Register page reachable at /login
// Allows you to login to access the rest of the service
func getRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if len(username) == 0 || len(password) == 0 {
			log.Fatal("Job form data is invalid")
		}

		if status, err := auth.Register(username, password); status == internal.RegisterSuccess && err == nil {
			// Account exists and the password is valid
			log.Println(internal.FormatRegisterStatus(status))
		} else {
			// Account is not valid
			log.Println(internal.FormatRegisterStatus(status))
			log.Println("Account could not be created for whatever reason")
		}
		return
	}
	tpl["register"].ExecuteTemplate(w, "base", nil)
}

// Job adding page reachable at /job/add
// Allows you to add a job
func getJobAdd(w http.ResponseWriter, r *http.Request) {
	has, err := auth.HasPermission(r)
	if err != nil {
		log.Println(err)
	}

	if !has {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

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
	tpl["add-job"].ExecuteTemplate(w, "base", nil)
}

func getPeers(w http.ResponseWriter, r *http.Request) {
	has, err := auth.HasPermission(r)
	if err != nil {
		log.Println(err)
	}
	tpl["peers"].ExecuteTemplate(w, "base", nil)
}

// An endpoint that rings the bell. Reachable at /ring?id=<job-id>
func getRing(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println(err)
	}

	id, err := strconv.Atoi(params["id"][0])
	if err != nil {
		log.Println(err)
	}

	internal.Ring(id)
}
