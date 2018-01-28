package handlers

import (
	"errors"
	"net/http"

	"github.com/fibreactive/chat/models"
	"github.com/gorilla/sessions"
)

//os.GetEnv(session_key)
var store = sessions.NewCookieStore([]byte("hello world"))

type UserHandler struct {
	models.UserService
}

func (uh *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
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
		session.Values["id"] = user.ID
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
	session.Values["id"] = ""
	if err := session.Save(r, w); err != nil {
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//func (uh *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}