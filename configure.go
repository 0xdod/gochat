package main

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/fibreactive/chat/handlers"
	"github.com/fibreactive/chat/models"

	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "damilola"
	password = "passme123"
	dbname   = "chat_golang"
)

var rh *handlers.RoomHandler
var uh *handlers.UserHandler
var th *handlers.TemplateHandler

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := GetDB(psqlInfo)
	if err != nil {
		panic(err)
	}
	ug := models.NewUserGorm(db)
	rg := models.NewRoomGorm(db)
	uh = &handlers.UserHandler{ug}
	rh = &handlers.RoomHandler{rg, ug, nil}
	rh.PopulateRooms()
	th = &handlers.TemplateHandler{ug, rg}
	db.AutoMigrate(&models.UserModel{}, &models.RoomModel{})
}

func GetDB(connInfo string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", connInfo)
	if err != nil {
		return nil, err
	}
	return db, err
}

func MapRoutes(n *negroni.Negroni) {
	r := mux.NewRouter()
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	guestRouter := r.NewRoute().Subrouter()
	guestRouter.HandleFunc("/login", uh.Login)
	guestRouter.HandleFunc("/signup", uh.Signup)

	authRouter := r.NewRoute().Subrouter()
	authRouter.Use(uh.MustAuth)
	authRouter.Handle("/upload", th.HandlePage("upload.html"))
	authRouter.Handle("/chat/{link}", th.HandlePage("chat.html"))
	authRouter.HandleFunc("/r/{id}", rh.Chat)
	authRouter.HandleFunc("/r/leave/{id}", rh.Leave).Methods("POST")
	authRouter.HandleFunc("/r", rh.ListOrCreate)
	authRouter.HandleFunc("/logout", uh.Logout)
	authRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/r", http.StatusSeeOther)
	})
	n.UseHandler(r)
}
