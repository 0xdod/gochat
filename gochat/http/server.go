package http

import (
	"embed"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"time"

	"github.com/0xdod/gochat/gochat"
	"github.com/0xdod/gochat/gochat/gorm"
	"github.com/0xdod/gochat/gochat/http/templates"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

//go:embed public
var publicFiles embed.FS

// Server wraps net/http server
type Server struct {
	server *http.Server
	router http.Handler
	*services
}

type services struct {
	user gochat.UserService
	room gochat.RoomService
}

// NewServer creates a new server.
func NewServer() *Server {
	server := &http.Server{
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	s := &Server{server: server}
	setupRoutes(s)
	return s
}

func (s *Server) CreateServices(db *gorm.DB) {
	s.services = &services{
		user: gorm.NewUserService(db),
	}
}

// LoadTemplates loads parsed templates into memory
func (s *Server) LoadTemplates() {
	templates.ParseTemplates()
}

func (s *Server) render(w io.Writer, name string, data interface{}) {
	templates.Render(w, name, data)
}

// Run the server
func (s *Server) Run(port string) {
	s.server.Addr = port
	log.Fatal(s.server.ListenAndServe())
}

func setupRoutes(s *Server) {
	r := mux.NewRouter().StrictSlash(true)
	n := negroni.Classic()
	r.PathPrefix("/public/").Handler(http.FileServer(http.FS(publicFiles)))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`Hello world!`))
	})
	r.HandleFunc("/chat", s.chat)
	r.HandleFunc("/ws", s.handleWS)
	r.HandleFunc("/login", s.login)
	r.HandleFunc("/signup", s.register)
	n.UseHandler(r)
	s.server.Handler = n
}

func dirname() string {
	if _, file, _, ok := runtime.Caller(0); ok {
		return filepath.Dir(file)
	}
	return "."
}
