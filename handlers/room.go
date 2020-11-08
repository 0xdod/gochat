package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/fibreactive/chat/models"

	"github.com/fibreactive/chat/chat"

	"github.com/gorilla/websocket"
)

//To get a web socket connection
var upgrader = &websocket.Upgrader{ReadBufferSize: chat.SocketBufferSize, WriteBufferSize: chat.SocketBufferSize}

type RoomHandler struct {
	rooms map[string]*chat.Room
	*chat.Room
}

func (rh *RoomHandler) handleLeaveRoom(w http.ResponseWriter, req *http.Request) {
	//find client
	user, ok := req.Context().Value("user").(*models.UserModel)
	if !ok {
		http.Error(w, errors.New("Cannot find client!").Error(), http.StatusForbidden)
		return
	}
	client := rh.Room.FindClient(user.ID)
	if client == nil {
		http.Error(w, errors.New("Cannot find client!").Error(), http.StatusForbidden)
		return
	}
	rh.Room.RemoveClient(client)
	w.WriteHeader(http.StatusOK)
	return
}

func (rh *RoomHandler) handleChat(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	user, ok := req.Context().Value("user").(*models.UserModel)
	if !ok {
		http.Error(w, errors.New("Cannot find client!").Error(), http.StatusForbidden)
		return
	}
	client := chat.NewClient(socket, user)
	defer rh.Room.RemoveClient(client)
	rh.Room.AddClient(client)
	go client.Write()
	client.Read()
}

func (rh *RoomHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/room":
		rh.handleChat(w, req)
	case "/leave":
		if req.Method == "POST" {
			rh.handleLeaveRoom(w, req)
		}
	}
}
