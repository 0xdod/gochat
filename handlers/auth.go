package handlers

import (
	"errors"
	"net/http"

	"github.com/fibreactive/chat/models"
	"github.com/gorilla/schema"

	"github.com/stretchr/objx"
)

var decoder = schema.NewDecoder()

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type SignupForm struct {
	Firstname string `schema:"firstname"`
	Lastname  string `schema:"lastname"`
	Nickname  string `schema:"nickname"`
	Email     string `schema:"email"`
	Password  string `schema:"password"`
	Password2 string `schema:"password2"`
}

func parseForm(r *http.Request, form interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	// r.PostForm is a map of our POST form values
	err = decoder.Decode(form, r.PostForm)
	if err != nil {
		return err
	}
	return nil
}

type authHandler struct {
	next http.Handler
}

func MustAuth(h http.Handler) http.Handler {
	return &authHandler{next: h}
}

func (auth *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err == http.ErrNoCookie || cookie.Value == "" {
		// not authenticated should login
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusSeeOther)
		return
	}
	if err != nil {
		//some other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// success - call the next handler
	auth.next.ServeHTTP(w, r)
}

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
	authCookieValue := objx.New(map[string]interface{}{
		"userid": user.Email,
		"name":   user.Nickname,
	}).MustBase64()

	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/"})
	w.Header().Set("Location", "/chat")
	w.WriteHeader(http.StatusSeeOther)
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
	w.Header().Set("Location", "/chat")
	w.WriteHeader(http.StatusSeeOther)
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
