package http

import (
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// Server wraps net/http server
type Server struct {
	server *http.Server
	router http.Handler
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

// Run the server
func (s *Server) Run(port string) {
	s.server.Addr = port
	log.Fatal(s.server.ListenAndServe())
}

func setupRoutes(s *Server) {
	r := mux.NewRouter().StrictSlash(true)
	n := negroni.Classic()
	dir, _ := filepath.Abs(filepath.Join(dirname(), "public"))
	r.PathPrefix("/public/").
		Handler(http.StripPrefix("/public/", http.FileServer(http.Dir(dir))))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`Hello world!`))
	})
	r.HandleFunc("/chat", index)
	r.HandleFunc("/ws", handleWS)
	n.UseHandler(r)
	s.server.Handler = n
}

func dirname() string {
	if _, file, _, ok := runtime.Caller(0); ok {
		return filepath.Dir(file)
	}
	return "."
}
