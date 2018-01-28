package http

import (
	"bytes"
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
		room: store.NewRoomStore(db),
		user: store.NewUserStore(db),
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
	buf := new(bytes.Buffer)
	err := templates.Render(buf, name, data)
	if err != nil {
		s.serverError(w, err)
		return
	}
	buf.WriteTo(w)
}

// Run the server
func (s *Server) Run(port string) {
	s.server.Addr = port
	s.ErrorLog.Fatal(s.server.ListenAndServe())
}

func setupRoutes(s *Server) {
	r := mux.NewRouter().StrictSlash(true)

	n := negroni.Classic()

	n.Use(negroni.HandlerFunc(SecurityMiddleware))
	n.Use(negroni.HandlerFunc(SessionMiddleware))
	n.Use(negroni.HandlerFunc(FlashMiddleware))

	auth := negroni.HandlerFunc(s.AuthMiddleware)

	r.PathPrefix("/public/").Handler(http.FileServer(http.FS(publicFiles)))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`Hello world!`))
	})
	r.Handle("/chat", adaptFunc(s.chat, auth))
	r.HandleFunc("/ws", s.handleWS)
	r.HandleFunc("/login", s.login)
	r.HandleFunc("/signup", s.register)
	r.Handle("/room/create", adaptFunc(s.createRoom, auth)).Methods("POST")
	r.Handle("/rooms", adaptFunc(s.listRooms, auth))

	n.UseHandler(r)

	s.server.Handler = n
}

func adaptNegroni(h http.Handler, nh ...negroni.Handler) http.Handler {
	n := negroni.New(nh...)
	return n.With(negroni.Wrap(h))
}

func adaptFunc(h http.HandlerFunc, nh ...negroni.Handler) http.Handler {
	n := negroni.New(nh...)
	return n.With(negroni.WrapFunc(h))

}
