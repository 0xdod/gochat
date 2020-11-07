package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fibreactive/chat/chat"
	"github.com/fibreactive/chat/handlers"
	"github.com/fibreactive/chat/models"

	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "damilola"
	password = "passme123"
	dbname   = "chat_golang"
)

var room = chat.CreateRoom("default")
var roomHandler = &handlers.RoomHandler{Room: room}
var userHandler *handlers.UserHandler

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	ug, err := models.NewUserGorm(psqlInfo)
	if err != nil {
		panic(err)
	}
	ug.AutoMigrate()
	userHandler = &handlers.UserHandler{ug}
}

func main() {
	r := mux.NewRouter()
	MapRoutes(r)
	go room.Run()
	log.Println("Starting server on port 9000")
	http.ListenAndServe(":9000", r)
}
