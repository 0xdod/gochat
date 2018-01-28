package templates

import (
	"errors"
	"html/template"
	"io"
	"log"
	"path/filepath"
	"sync"

	"io/fs"

	"embed"
)

//go:embed layouts/*.html pages/*.html partials/*.html
var templatesFS embed.FS

type templateStore struct {
	once      sync.Once
	templates map[string]*template.Template
	isParsed  bool
}

var globalTemplates = &templateStore{}

// ParseTemplates parses the templates in a given directory.
func (ts *templateStore) ParseTemplates() {
	ts.once.Do(func() {
		templates := make(map[string]*template.Template)
		layouts, err := fs.Glob(templatesFS, "layouts/*.html")
		partials, err := fs.Glob(templatesFS, "partials/*.html")
		pages, err := fs.Glob(templatesFS, "pages/*.html")

		if err != nil {
			log.Fatal(err)
		}

		for _, layout := range layouts {
			for _, page := range pages {
				files := []string{layout, page}
				files = append(files, partials...)
				temp := template.Must(template.ParseFS(templatesFS, files...))
				templates[filepath.Base(page)] = temp
			}
		}
		ts.templates = templates
		ts.isParsed = true
	})
}

// Render executes a template to an io.Writer.
func (ts *templateStore) Render(out io.Writer, name string, data interface{}) error {
	if !ts.isParsed {
		return errors.New("templates not parsed")
	}

	return ts.templates[name].ExecuteTemplate(out, "base", data)
}

// ParseTemplates parses the files in a given directory
// to the global templateStore.
func ParseTemplates() {
	globalTemplates.ParseTemplates()
}

// Render executes a template from the global templateStore.
func Render(out io.Writer, templateName string, contextData interface{}) error {
	return globalTemplates.Render(out, templateName, contextData)
}
