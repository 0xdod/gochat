package http

import (
	"embed"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/0xdod/gochat/gochat"
	"github.com/0xdod/gochat/gochat/http/templates"
	"github.com/0xdod/gochat/gochat/store"
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
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

type services struct {
	user gochat.UserService
	room gochat.RoomService
}

type templateData map[string]interface{}

// NewServer creates a new server.
func NewServer() *Server {
	server := &http.Server{
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	s := &Server{server: server}
	s.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	s.ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	setupRoutes(s)
	return s
}

func (s *Server) CreateServices(db *store.DB) {
	s.services = &services{
		user: store.NewGormStore(db),
	}
}

// LoadTemplates loads parsed templatesates into memory
func (s *Server) LoadTemplates() {
	templates.ParseTemplates()
}

func (s *Server) render(w http.ResponseWriter, r *http.Request, name string, data templateData) {
	if data == nil {
		data = templateData{}
	}
	data["request"] = r
	data["user"] = UserFromContext(r.Context())
	templates.Render(w, name, data)
}

// Run the server
func (s *Server) Run(port string) {
	s.server.Addr = port
	s.ErrorLog.Fatal(s.server.ListenAndServe())
}

func setupRoutes(s *Server) {
	r := mux.NewRouter().StrictSlash(true)

	n := negroni.Classic()

	nn := negroni.New(
		negroni.HandlerFunc(s.AuthMiddleware),
		negroni.Wrap(http.HandlerFunc(s.roomList)))

	r.PathPrefix("/public/").Handler(http.FileServer(http.FS(publicFiles)))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`Hello world!`))
	})
	r.HandleFunc("/chat", s.chat)
	r.HandleFunc("/ws", s.handleWS)
	r.HandleFunc("/login", s.login)
	r.HandleFunc("/signup", s.register)
	r.HandleFunc("/room", s.createRoom)
	r.Handle("/rooms", nn)

	n.Use(negroni.HandlerFunc(SessionMiddleware))
	n.Use(negroni.HandlerFunc(FlashMiddleware))
	n.UseHandler(r)

	s.server.Handler = n
}
