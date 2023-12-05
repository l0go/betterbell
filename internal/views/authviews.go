package views

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"codeberg.org/logo/betterbell/internal"
)

type AuthViews struct {
	State            internal.AuthState
	LoginTemplate    *template.Template
	RegisterTemplate *template.Template
}

func CreateAuth(db *sql.DB) AuthViews {
	var v AuthViews

	// Create the authentication state
	v.State = internal.AuthState{
		DB:       db,
		Sessions: make(map[[16]byte]struct{}),
	}

	// Create templates
	v.LoginTemplate = template.Must(template.ParseFiles("templates/login.html", "templates/base.html"))
	v.RegisterTemplate = template.Must(template.ParseFiles("templates/register.html", "templates/base.html"))

	return v
}

// Allows you to login to access the rest of the service
func (v AuthViews) GetLogin(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println(err)
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if len(username) == 0 || len(password) == 0 {
			log.Fatal("Job form data is invalid")
		}

		if status, err := v.State.Check(username, password); status == internal.LoginSuccess && err == nil {
			// Account exists and the password is valid
			session, err := v.State.GrantSession(r)
			if err != nil {
				log.Println(err)
				return
			}
			http.SetCookie(w, &session)

			if _, ok := params["redirect"]; ok {
				http.Redirect(w, r, params["redirect"][0], http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		} else {
			// Account is not valid
			http.Redirect(w, r, fmt.Sprintf("/login?error=%d", status), http.StatusSeeOther)
		}
		return
	}

	if _, ok := params["error"]; ok {
		errorStatus, _ := strconv.Atoi(params["error"][0])
		v.LoginTemplate.ExecuteTemplate(w, "base", internal.FormatLoginStatus(internal.LoginStatus(errorStatus)))
		return
	}

	v.LoginTemplate.ExecuteTemplate(w, "base", nil)
}

// Register page reachable at /register
// Allows you to login to access the rest of the service
func (v AuthViews) GetRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if len(username) == 0 || len(password) == 0 {
			log.Fatal("Job form data is invalid")
		}

		status, err := v.State.Register(username, password)
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
		v.RegisterTemplate.ExecuteTemplate(w, "base", internal.FormatRegisterStatus(internal.RegisterStatus(errorStatus)))
		return
	}

	v.RegisterTemplate.ExecuteTemplate(w, "base", nil)
}

func (v AuthViews) CheckPermission(w http.ResponseWriter, r *http.Request) {
	has, err := v.State.HasPermission(r)
	if err != nil {
		log.Println(err)
	}

	if !has {
		path := r.URL.Path
		http.Redirect(w, r, fmt.Sprintf("/login?redirect=%s", path), http.StatusSeeOther)
		return
	}
}
