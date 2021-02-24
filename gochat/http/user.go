package http

import "net/http"

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	s.render(w, "login.html", r)
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	s.render(w, "signup.html", r)
}
