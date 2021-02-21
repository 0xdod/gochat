package templates

import (
	"html/template"
	"io"
	"log"
	"path/filepath"
	"runtime"
	"sync"
)

// type TemplateConfig struct {
// 	LayoutsDir  string
// 	PartialsDir string
// }

// Template is used to hold multiple templates
// to allow sharing of layout template
type templateStore struct {
	once      sync.Once
	templates map[string]*template.Template
	isParsed  bool
}

var globalTemplates *templateStore = &templateStore{}

func (ts *templateStore) Templates() map[string]*template.Template {
	return ts.templates
}

func (ts *templateStore) ParseTemplates(templatesDir string) {
	ts.once.Do(func() {
		templates := make(map[string]*template.Template)
		layoutsDir := filepath.Join(templatesDir, "layouts")
		partialsDir := filepath.Join(templatesDir, "partials")
		layouts, err := filepath.Glob(layoutsDir + "/*.html")
		partials, err := filepath.Glob(partialsDir + "/*.html")
		pages, err := filepath.Glob(templatesDir + "/pages/*.html")
		if err != nil {
			log.Fatal(err)
		}

		for _, layout := range layouts {
			for _, page := range pages {
				files := []string{layout, page}
				files = append(files, partials...)
				temp := template.Must(template.ParseFiles(files...))
				templates[filepath.Base(page)] = temp
			}
		}
		ts.templates = templates
		ts.isParsed = true
	})
}

func Render(out io.Writer, templateName string, contextData interface{}) {
	if !globalTemplates.isParsed {
		globalTemplates.ParseTemplates(dirname())
	}
	t := globalTemplates.Templates()
	t[templateName].ExecuteTemplate(out, "base", contextData)
}

func dirname() string {
	if _, file, _, ok := runtime.Caller(0); ok {
		return filepath.Dir(file)
	}
	return "."
}
