package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/0xdod/gochat/handlers"
	"github.com/0xdod/gochat/models"

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
var mh *handlers.MessageHandler

func init() {
	localPsql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	psqlInfo := os.Getenv("DATABASE_URL")

	if psqlInfo == "" {
		psqlInfo = localPsql
	}

	db, err := GetDB(psqlInfo)
	if err != nil {
		panic(err)
	}
	ug := models.NewUserGorm(db)
	rg := models.NewRoomGorm(db)
	mg := models.NewMessageGorm(db)
	uh = &handlers.UserHandler{ug, rg}
	rh = &handlers.RoomHandler{rg, ug, nil}
	th = &handlers.TemplateHandler{ug, rg}
	mh = &handlers.MessageHandler{mg, rg, ug}
	db.AutoMigrate(&models.User{}, &models.Room{}, &models.Message{})
	rh.PopulateRooms()
}

func GetDB(connInfo string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", connInfo)
	if err != nil {
		return nil, err
	}
	return db, err
}

func MapRoutes(n *negroni.Negroni) {
	r := mux.NewRouter().StrictSlash(true)
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	guestRouter := r.NewRoute().Subrouter()
	guestRouter.HandleFunc("/login", uh.Login)
	guestRouter.HandleFunc("/signup", uh.Signup)

	authRouter := r.NewRoute().Subrouter()
	authRouter.Use(uh.MustAuth)
	authRouter.Handle("/upload", th.HandlePage("upload.html"))
	authRouter.Handle("/chat/{link}", th.HandlePage("chat.html"))
	authRouter.HandleFunc("/r/{id}", rh.ServeWs)
	authRouter.HandleFunc("/r/leave/{id}", rh.Leave).Methods("POST")
	authRouter.HandleFunc("/r", rh.ListOrCreate)
	authRouter.HandleFunc("/logout", uh.Logout)
	authRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/r", http.StatusSeeOther)
	})
	authRouter.HandleFunc("/messages", mh.Create).Methods("POST")
	authRouter.HandleFunc("/messages", mh.List).Methods("GET")
	authRouter.HandleFunc("/u/{id}", uh.ProfileDetail)
	n.UseHandler(r)
}
