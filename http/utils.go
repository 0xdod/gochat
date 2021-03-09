package http

import (
	"bytes"
	"encoding/gob"
	"log"
	"net/http"
	"strings"

	"github.com/0xdod/gochat"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

var schemaDecoder = schema.NewDecoder()
var validate = validator.New()
var sessionStore = sessions.NewCookieStore([]byte("some-really-deep-secret."))

type FlashMessage struct {
	Type    string
	Message string
	Style   string
}

func NewFlash(key, message string) *FlashMessage {
	var style string
	switch key {
	case "success":
		style = "success"
	case "error":
		style = "danger"
	case "info":
		style = "info"
	default:
		style = "warning"

	}
	return &FlashMessage{
		Type:    key,
		Message: message,
		Style:   style,
	}
}

func (fm *FlashMessage) Serialize() string {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(fm)
	if err != nil {
		log.Println(err)
	}
	return buf.String()
}

func DeserializeFlashMessage(str string) FlashMessage {
	var fm FlashMessage
	dec := gob.NewDecoder(strings.NewReader(str))
	err := dec.Decode(&fm)
	if err != nil {
		log.Println(err)
	}
	return fm
}

func validateStruct(s interface{}) error {
	return validate.Struct(s)
}

func addFlash(w http.ResponseWriter, r *http.Request, key, value string) {
	session, _ := sessionStore.Get(r, "flash")
	f := NewFlash(key, value)
	session.AddFlash(f.Serialize(), "flash")
	session.Save(r, w)
}

func loginSession(w http.ResponseWriter, r *http.Request, user *gochat.User) {
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
