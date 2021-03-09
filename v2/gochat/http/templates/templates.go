package templates

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"path/filepath"
	"sync"

	"io/fs"

	"embed"
)

type Options struct {
	fs.FS
	MainLayout string
	Layouts    string
	Partials   string
	Pages      string
	Ext        string
	Funcs      template.FuncMap
	Directory  string
}

func New(opts ...Options) *templateStore {
	var options Options
	if len(opts) > 0 {
		options = opts[0]
	} else {
		options = defaultOptions
	}
	tmpl := &templateStore{
		options: options,
	}
	tmpl.Parse()
	return tmpl
}

//go:embed layouts/*.html pages/*.html partials/*.html
var templatesFS embed.FS

type templateStore struct {
	once      sync.Once
	templates map[string]*template.Template
	parsed    bool
	options   Options
}

var defaultOptions = Options{
	FS:         templatesFS,
	MainLayout: "base",
	Layouts:    "layouts",
	Partials:   "partials",
	Ext:        ".html",
	Pages:      "pages",
}

// Parse parses the templates in a given directory.
func (ts *templateStore) Parse() {
	ts.once.Do(func() {
		opts := ts.options
		layoutsGlob := fmt.Sprintf("%s/*%s", opts.Layouts, opts.Ext)
		partialsGlob := fmt.Sprintf("%s/*%s", opts.Partials, opts.Ext)
		pagesGlob := fmt.Sprintf("%s/*%s", opts.Pages, opts.Ext)
		layouts, err := fs.Glob(opts.FS, layoutsGlob)
		partials, err := fs.Glob(opts.FS, partialsGlob)
		pages, err := fs.Glob(opts.FS, pagesGlob)

		if err != nil {
			log.Fatal(err)
		}
		templates := make(map[string]*template.Template)

		for _, page := range pages {
			name := filepath.Base(page)
			files := make([]string, len(layouts))
			copy(files, layouts)
			files = append(files, page)
			files = append(files, partials...)
			temp := template.Must(template.ParseFS(opts.FS, files...))
			templates[name] = temp
		}
		ts.templates = templates
		ts.parsed = true
	})
}

// Render executes a template to an io.Writer.
func (ts *templateStore) Render(out io.Writer, name string, data interface{}) error {
	if !ts.parsed {
		return errors.New("templates not parsed")
	}
	mainLayout := ts.options.MainLayout
	return ts.templates[name].ExecuteTemplate(out, mainLayout, data)
}

type Renderer interface {
	Render(out io.Writer, name string, data interface{}) error
}
