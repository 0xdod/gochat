package handlers

import (
	"net/http"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

type CreateMessageJSON struct {
	Message string `json:"message,omitempty"`
	RoomID  int    `json:"roomID,omitempty"`
	UserID  int    `json:"userID,omitempty"`
}

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

type CreateRoomForm struct {
	Name        string `schema:"name"`
	Description string `schema:"description"`
	AvatarURL   string `schema:"omitempty"`
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
