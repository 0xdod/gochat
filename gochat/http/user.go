package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/0xdod/gochat/gochat"
	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

var schemaDecoder = schema.NewDecoder()
var validate = validator.New()
var sessionStore = sessions.NewCookieStore([]byte("some-really-deep-secret."))

type userSignUpForm struct {
	Name      string `json:"name" schema:"name" validate:"required"`
	Username  string `json:"username" schema:"username" validate:"required"`
	Email     string `json:"email" schema:"email" validate:"required,email"`
	Password  string `json:"password" schema:"password" validate:"required,min=8"`
	Password2 string `json:"password2" schema:"password2" validate:"required,min=8"`
}

type userLoginForm struct {
	Email    string `json:"email" schema:"email" validate:"required,email"`
	Password string `json:"password" schema:"password" validate:"required,min=8"`
}

func (form *userSignUpForm) create() *gochat.User {
	user := &gochat.User{
		Name:     form.Name,
		Username: form.Username,
		Email:    form.Email,
	}
	_ = user.SetPassword(form.Password)
	return user
}

func validateStruct(s interface{}) error {
	return validate.Struct(s)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		form := new(userLoginForm)
		_ = parseForm(r, form)
		if err := validateStruct(form); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user := s.services.user.Authenticate(context.Background(), form.Email, form.Password)
		if user == nil {
			AddFlash(w, r, "Cannot find user with such details.")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		AddFlash(w, r, "Welcome!")
		LoginSession(w, r, user)
		http.Redirect(w, r, "/rooms", http.StatusSeeOther)
		return
		//login and start session.
		// redirect to chat room list.
	}
	s.render(w, "login.html", r)
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		form := new(userSignUpForm)
		_ = parseForm(r, form)
		if err := validateStruct(form); err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := s.services.user.CreateUser(context.Background(), form.create())
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		AddFlash(w, r, "User created.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	s.render(w, "signup.html", r)
}

func FlashMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session, _ := sessionStore.Get(r, "flash")
	fm := session.Flashes("flash")
	fmt.Println(fm)
	session.Save(r, w)
	r = r.WithContext(context.WithValue(r.Context(), "messages", fm))
	next(w, r)

}

func AddFlash(w http.ResponseWriter, r *http.Request, value interface{}) {

	session, _ := sessionStore.Get(r, "flash")
	session.AddFlash(value, "flash")
	session.Save(r, w)

}

func LoginSession(w http.ResponseWriter, r *http.Request, user *gochat.User) {
	session, _ := sessionStore.Get(r, "session.id")
	session.Values["is_auth"] = true
	session.Values["user_id"] = user.ID
	session.Save(r, w)
}

func parseForm(r *http.Request, form interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return schemaDecoder.Decode(form, r.PostForm)
}
