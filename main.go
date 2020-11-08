package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fibreactive/chat/chat"
	"github.com/fibreactive/chat/handlers"
	"github.com/fibreactive/chat/models"
	"github.com/urfave/negroni"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "damilola"
	password = "passme123"
	dbname   = "chat_golang"
)

var room = chat.NewRoom("default")
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
	n := negroni.Classic()
	MapRoutes(n)
	go room.Run()
	log.Println("Starting server on port 9000")
	http.ListenAndServe(":9000", n)
}
