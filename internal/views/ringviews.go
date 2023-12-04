package views

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"codeberg.org/logo/betterbell/internal"
)

type RingViews struct {
	Peers PeerViews
}

// An endpoint that rings the bell. Reachable at /ring/?id=<job-id>
func (v RingViews) GetRing(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println(err)
	}

	id, err := strconv.Atoi(params["id"][0])
	if err != nil {
		log.Println(err)
	}

	secret := params["secret"][0]

	valid, err := v.Peers.State.Check(secret)
	if err != nil {
		log.Println(err)
	}

	if valid {
		internal.Ring(id, &v.Peers.State)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
