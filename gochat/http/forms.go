package http

import (
	"strings"

	"github.com/0xdod/gochat/gochat"
)

type form interface {
	validate() error
}

type baseForm struct {
	Errors map[string]string
}

type userSignUpForm struct {
	Name      string            `json:"name" schema:"name" validate:"required"`
	Username  string            `json:"username" schema:"username" validate:"required"`
	Email     string            `json:"email" schema:"email" validate:"required,email"`
	Password  string            `json:"password" schema:"password" validate:"required,min=8"`
	Password2 string            `json:"password2" schema:"password2" validate:"required,min=8"`
	Errors    map[string]string `json:"-" schema:"-" validate:"-"`
}

type userLoginForm struct {
	Email    string            `json:"email" schema:"email" validate:"required,email"`
	Password string            `json:"password" schema:"password" validate:"required,min=8"`
	Errors   map[string]string `json:"-" schema:"-" validate:"-"`
}

func (form *userSignUpForm) create() *gochat.User {
	email := strings.ToLower(form.Email)
	user := &gochat.User{
		Name:     form.Name,
		Username: form.Username,
		Email:    email,
	}
	_ = user.SetPassword(form.Password)
	return user
}

func (form *userSignUpForm) validate() bool {
	if form.Password == form.Password2 {
		form.Errors["password"] = "Passwords do not match."
	}
	return false
}
