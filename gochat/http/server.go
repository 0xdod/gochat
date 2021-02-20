package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server wraps net/http server
type Server struct {
	server *http.Server
	router http.Handler
}

// NewServer creates a new server.
func NewServer() *Server {
	r := mux.NewRouter()
	s = &http.Server{
		Handler: r,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return &Server{s, r}
}

// Run the servr
func (s *Server) Run(port string) {
	s.server.Addr = port
	log.Fatal(s.server.ListenAndServe())
}

// Router return mux
func (s *Server) Router() *mux.Router {
	return s.router.(*mux.Router)
}
