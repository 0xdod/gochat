package main

import (
	"net/http"
	"time"

	h "github.com/fibreactive/chat/handlers"

	"github.com/gorilla/mux"
)

func MapRoutes(r *mux.Router) {
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	r.Handle("/chat", h.MustAuth(h.HandlePage("chat.html")))
	r.Handle("/signup", h.HandlePage("signup.html")).Methods("GET")
	r.Handle("/login", h.HandlePage("login.html")).Methods("GET")
	r.Handle("/login", userHandler).Methods("POST")
	r.Handle("/signup", userHandler).Methods("POST")
	r.Handle("/upload", h.MustAuth(h.HandlePage("upload.html"))).Methods("GET")
	r.Handle("/room", h.MustAuth(room))
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    "auth",
			Expires: time.Now(),
			Value:   "",
			Path:    "/",
			MaxAge:  -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
}
