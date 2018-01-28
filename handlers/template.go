package handlers

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/fibreactive/chat/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/objx"
)

type TemplateHandler struct {
	models.UserService
	models.RoomService
}

type templateStore struct {
	once      sync.Once
	templates map[string]*template.Template
}

func (ts *templateStore) renderTemplate(out io.Writer, filename string, data interface{}) {
	ts.once.Do(ts.parseTemplates)
	ts.templates[filename].ExecuteTemplate(out, "base", data)
}

func (ts *templateStore) renderOne(out io.Writer, filename string, data interface{}) {
	ts.once.Do(ts.parseTemplates)
	ts.templates["chat.html"].ExecuteTemplate(out, filename, data)
}

func (ts *templateStore) parseTemplates() {
	// lazy load and compile template only once
	ts.templates = make(map[string]*template.Template)
	templatesDir := "templates"
	layoutsDir := filepath.Join(templatesDir, "layouts")
	partialsDir := filepath.Join(templatesDir, "partials")
	layouts, err := filepath.Glob(layoutsDir + "/*.html")
	partials, err := filepath.Glob(partialsDir + "/*.html")
	pages, err := filepath.Glob(templatesDir + "/*.html")
	if err != nil {
		log.Fatal(err)
	}
	for _, layout := range layouts {
		for _, page := range pages {
			files := []string{layout, page}
			files = append(files, partials...)
			temp := template.Must(template.ParseFiles(files...))
			ts.templates[filepath.Base(page)] = temp
		}
	}
}

var ts = &templateStore{}

func render(out io.Writer, templateName string, data interface{}) {
	ts.renderTemplate(out, templateName, data)
}

func renderOne(out io.Writer, filename string, data interface{}) {
	ts.renderOne(out, filename, data)
}

func (t *TemplateHandler) HandlePage(templateName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			m := objx.MSI()
			m.Set("Host", r.Host)
			user, ok := Get(r, "user").(*models.User)
			if ok {
				m.Set("user", user)
			}
			vars := mux.Vars(r)
			if link, exists := vars["link"]; exists {
				userRooms := t.UserService.GetRooms(user)
				room := t.RoomService.FindByLink(link)
				t.RoomService.AddParticipant(room, user)
				present := false
				for _, r := range userRooms {
					if r.ID == room.ID {
						present = true
						break
					}
				}
				session, _ := store.Get(r, "session.id")
				session.Values["room_id"] = room.ID
				session.Values["present"] = present
				session.Save(r, w)
				m.Set("room", room)
			}
			m.Set("MyRooms", t.UserService.GetRooms(user))
			render(w, templateName, m)
		}
	})
}
