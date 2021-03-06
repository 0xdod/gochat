package handlers

import (
	"errors"
	"net/http"
	"os"

	"github.com/0xdod/gochat/models"
	"github.com/gorilla/sessions"
	"github.com/stretchr/objx"
)

var sessionKey = getEnv("SESSION_KEY", "!(#IRJfqkfjwKEENRLDNnA")
var store = sessions.NewCookieStore([]byte(sessionKey))

type UserHandler struct {
	models.UserService
	models.RoomService
}

func getEnv(key, def string) string {
	s := os.Getenv(key)
	if s == "" {
		return def
	}
	return s
}

func (uh *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		form := SignupForm{}
		if err := parseForm(r, &form); err != nil {
			panic(err)
		}
		user := &models.User{
			Firstname: form.Firstname,
			Lastname:  form.Lastname,
			Username:  form.Username,
			Email:     form.Email,
			Password:  form.Password,
		}
		if err := uh.UserService.Create(user); err != nil {
			panic(err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	render(w, "signup.html", nil)
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		form := LoginForm{}
		if err := parseForm(r, &form); err != nil {
			panic(err)
		}
		user := uh.UserService.Authenticate(form.Email, form.Password)
		if user == nil {
			http.Error(w, errors.New("Forbidden request").Error(), http.StatusForbidden)
			return
		}
		session, _ := store.Get(r, "session.id")
		session.Values["user_id"] = user.ID
		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	render(w, "login.html", nil)
}

func (*UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session.id")
	session.Values["user_id"] = ""
	if err := session.Save(r, w); err != nil {
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (uh *UserHandler) ProfileDetail(w http.ResponseWriter, r *http.Request) {
	user := Get(r, "user")
	data := objx.MSI()
	data.Set("user", user)
	render(w, "user_profile.html", data)
}

//func (uh *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
