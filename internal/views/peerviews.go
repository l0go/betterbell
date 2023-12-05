package views

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"codeberg.org/logo/betterbell/internal"
)

type PeerViews struct {
	State        internal.PeerState
	ListTemplate *template.Template
	Auth         *AuthViews
}

func CreatePeers(db *sql.DB, auth *AuthViews) PeerViews {
	var v PeerViews
	v.State = internal.PeerState{DB: db}
	v.ListTemplate = template.Must(template.ParseFiles("templates/peers.html", "templates/base.html"))
	v.Auth = auth
	return v
}

func (v PeerViews) GetPeers(w http.ResponseWriter, r *http.Request) {
	// Check if we have permission
	v.Auth.CheckPermission(w, r)

	if r.Method == http.MethodPost {
		endpoint := r.FormValue("endpoint")
		secret := r.FormValue("secret")

		if len(endpoint) == 0 || len(secret) == 0 {
			log.Fatal("Job form data is invalid")
		}

		if err := v.State.Add(endpoint, secret); err != nil {
			log.Fatal(err)
		}
	}

	v.ListTemplate.ExecuteTemplate(w, "base", v.State)
}

// An endpoint that removes a peer. Reachable at /peers/remove/?id=<peer-id>
func (v PeerViews) GetPeersRemove(w http.ResponseWriter, r *http.Request) {
	// Check if we have permission
	v.Auth.CheckPermission(w, r)

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println(err)
	}

	id, err := strconv.Atoi(params["id"][0])
	if err != nil {
		log.Println(err)
	}

	v.State.Remove(id)
	http.Redirect(w, r, "/peers/", http.StatusSeeOther)
}
