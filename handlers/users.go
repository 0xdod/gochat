package handlers

import (
	"errors"
	"net/http"
	"os"

	"github.com/fibreactive/chat/models"
	"github.com/gorilla/sessions"
)

// var key string

// func init() {
// 	key = os.Getenv("SESSION_KEY")
// 	if key == "" {
// 		os.Setenv("SESSION_KEY", string(securecookie.GenerateRandomKey(32)))
// 	}
// }

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var Store = store

type UserHandler struct {
	models.UserService
}

func (uh *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
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
	session.Values["id"] = user.ID
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

func (uh *UserHandler) handleSignup(w http.ResponseWriter, r *http.Request) {
	form := SignupForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := &models.UserModel{
		Firstname: form.Firstname,
		Lastname:  form.Lastname,
		Nickname:  form.Nickname,
		Email:     form.Email,
		Password:  form.Password,
	}
	if err := uh.UserService.Create(user); err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

func (uh *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		switch r.URL.Path {
		case "/login":
			uh.handleLogin(w, r)
		case "/signup":
			uh.handleSignup(w, r)
		default:
			http.Error(w, errors.New("Bad request").Error(), http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, errors.New("Not implemented").Error(), http.StatusInternalServerError)
		return
	}
}
