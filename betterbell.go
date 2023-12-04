package main

import (
	"fmt"
	"os"
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

var (
	tpl   map[string]*template.Template
	auth  internal.AuthState
	jobs  internal.JobsState
	peers internal.PeerState
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

	// Create the gocron scheduler
	scheduler := gocron.NewScheduler(time.Local)

	// Create the peer state
	peers = internal.PeerState{DB: db}

	// Create the jobs state
	jobs = internal.JobsState{
		Scheduler: scheduler,
		Peers:     &peers,
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
	http.HandleFunc("/peers/remove/", getPeersRemove)
	http.HandleFunc("/ring/", getRing)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err = http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Index reachable at /
// Contains a list of the jobs and actions that can be acted upon them
func getRoot(w http.ResponseWriter, r *http.Request) {
	// Check if we have permission
	checkPermission(w, r)

	if r.Method == http.MethodPost {
		action, id, _ := strings.Cut(r.FormValue("action"), "-")
		switch action {
		case "ring":
			intId, _ := strconv.Atoi(id)
			if err := internal.Ring(intId, &peers); err != nil {
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
	tpl["index"].ExecuteTemplate(w, "base", jobs.Get())
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
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			// Account is not valid
			log.Println(internal.FormatLoginStatus(status))
		}
		return
	}
	tpl["login"].ExecuteTemplate(w, "base", nil)
}

// Register page reachable at /register
// Allows you to login to access the rest of the service
func getRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if len(username) == 0 || len(password) == 0 {
			log.Fatal("Job form data is invalid")
		}

		status, err := auth.Register(username, password)
		if err != nil {
			log.Fatal(err)
		}

		if status == internal.RegisterSuccess {
			// Account exists and the password is valid
			log.Println(internal.FormatRegisterStatus(status))
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			// Account is not valid
			http.Redirect(w, r, fmt.Sprintf("/register?error=%d", status), http.StatusSeeOther)
		}
		return
	}

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println(err)
	}

	if _, ok := params["error"]; ok {
		errorStatus, _ := strconv.Atoi(params["error"][0])
		tpl["register"].ExecuteTemplate(w, "base", internal.FormatRegisterStatus(internal.RegisterStatus(errorStatus)))
		return
	}

	tpl["register"].ExecuteTemplate(w, "base", nil)
}

// Job adding page reachable at /job/add
// Allows you to add a job
func getJobAdd(w http.ResponseWriter, r *http.Request) {
	// Check if we have permission
	checkPermission(w, r)
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
	// Check if we have permission
	checkPermission(w, r)

	if r.Method == http.MethodPost {
		endpoint := r.FormValue("endpoint")
		secret := r.FormValue("secret")

		if len(endpoint) == 0 || len(secret) == 0 {
			log.Fatal("Job form data is invalid")
		}

		if err := peers.Add(endpoint, secret); err != nil {
			log.Fatal(err)
		}
	}

	tpl["peers"].ExecuteTemplate(w, "base", peers.Get())
}

// An endpoint that removes a peer. Reachable at /peers/remove/?id=<peer-id>
func getPeersRemove(w http.ResponseWriter, r *http.Request) {
	// Check if we have permission
	checkPermission(w, r)

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println(err)
	}

	id, err := strconv.Atoi(params["id"][0])
	if err != nil {
		log.Println(err)
	}

	peers.Remove(id)
	http.Redirect(w, r, "/peers/", http.StatusSeeOther)
}

// An endpoint that rings the bell. Reachable at /ring/?id=<job-id>
func getRing(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println(err)
	}

	id, err := strconv.Atoi(params["id"][0])
	if err != nil {
		log.Println(err)
	}

	secret := params["secret"][0]

	valid, err := peers.Check(secret)
	if err != nil {
		log.Println(err)
	}

	if valid {
		internal.Ring(id, &peers)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func checkPermission(w http.ResponseWriter, r *http.Request) {
	has, err := auth.HasPermission(r)
	if err != nil {
		log.Println(err)
	}

	if !has {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}
