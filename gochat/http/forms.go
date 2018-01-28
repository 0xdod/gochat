package http

import (
	"strings"
	"unicode/utf8"

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

type roomCreateForm struct {
	Name        string  `json:"name" schema:"name" validate:"required,max=8"`
	Description *string `json:"description,omitempty" schema:"description"`
	Icon        *string `json:"icon,omitempty" schema:"icon"`
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

func (form *roomCreateForm) validate() map[string]string {
	errors := make(map[string]string)

	// err := validateStruct(form)
	// if err != nil {
	// 	switch e := err.(type) {
	// 	case validator.InvalidValidationError:
	// 		return e
	// 	case validator.ValidationErrors:
	// 		fmt.Println(e)
	// 		return e
	// 	}
	// }

	if strings.TrimSpace(form.Name) == "" {
		errors["name"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Name) > 100 {
		errors["name"] = "This field is too long (maximum is 100 characters)"
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (form *roomCreateForm) create() *gochat.Room {
	return &gochat.Room{
		Name:        form.Name,
		Description: form.Description,
		Icon:        form.Icon,
	}
}
