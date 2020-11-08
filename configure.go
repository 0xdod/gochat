package main

import (
	"net/http"

	h "github.com/fibreactive/chat/handlers"
	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
)

func MapRoutes(n *negroni.Negroni) {
	r := mux.NewRouter()
	authRouter := r.PathPrefix("").Subrouter()
	authRouter.Use(userHandler.MustAuth)
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	authRouter.Handle("/chat", h.HandlePage("chat.html")).Methods("GET")
	authRouter.Handle("/upload", h.HandlePage("upload.html")).Methods("GET")
	r.Handle("/signup", h.HandlePage("signup.html")).Methods("GET")
	r.Handle("/login", h.HandlePage("login.html")).Methods("GET")
	r.Handle("/login", userHandler).Methods("POST")
	r.Handle("/signup", userHandler).Methods("POST")
	authRouter.Handle("/room", roomHandler)
	authRouter.Handle("/leave", roomHandler).Methods("POST")
	authRouter.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.Store.Get(r, "session.id")
		session.Values["id"] = ""
		if err := session.Save(r, w); err != nil {
			return
		}
		http.Redirect(w, r, "/chat", http.StatusSeeOther)
	})
	n.UseHandler(r)
}
