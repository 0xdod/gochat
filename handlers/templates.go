package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/objx"
)

func HandlePage(page string) http.Handler {
	return &templateHandler{filename: page}
}

type templateHandler struct {
	filename string
	once     sync.Once
	templ    map[string]*template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// load and compile template only once
	t.once.Do(func() {
		t.templ = make(map[string]*template.Template)
		//load main base layout dir
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
				t.templ[filepath.Base(page)] = temp
			}
		}
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if err, ok := r.Context().Value("error").(error); ok {
		data["Error"] = err
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ[t.filename].ExecuteTemplate(w, "base", data)
}
