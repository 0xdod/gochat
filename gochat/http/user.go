package http

import (
	"context"
	"net/http"

	"github.com/0xdod/gochat/gochat"
)

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		form := new(userLoginForm)
		_ = parseForm(r, form)
		if err := validateStruct(form); err != nil {
			//http.Error(w, err.Error(), http.StatusBadRequest)
			s.clientError(w, http.StatusBadRequest)
			return
		}
		user := s.services.user.Authenticate(context.Background(), form.Email, form.Password)
		if user == nil {
			addFlash(w, r, "error", "Cannot find user with such details.")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		addFlash(w, r, "success", "Welcome "+user.Username)
		loginSession(w, r, user)
		next := r.URL.Query().Get("next")
		if next == "" {
			next = "/rooms"
		}
		http.Redirect(w, r, next, http.StatusSeeOther)
		return
	}
	s.render(w, r, "login.html", nil)
}

func NewContextWithUser(ctx context.Context, user *gochat.User) context.Context {
	return context.WithValue(ctx, "user", user)
}

func UserFromContext(ctx context.Context) *gochat.User {
	user, ok := ctx.Value("user").(*gochat.User)
	if !ok {
		return nil
	}
	return user
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		form := new(userSignUpForm)
		_ = parseForm(r, form)
		if err := validateStruct(form); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := s.services.user.CreateUser(context.Background(), form.create())
		if err != nil {
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			if err == gochat.ECONFLICT {
				s.render(w, r, "signup.html", templateData{"error": "User with this email exists"})
			} else {
				s.serverError(w, err)
			}
			return
		}
		addFlash(w, r, "success", "User created.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	s.render(w, r, "signup.html", nil)
}
