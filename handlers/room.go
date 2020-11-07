package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/fibreactive/chat/chat"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

//To get a web socket connection
var upgrader = &websocket.Upgrader{ReadBufferSize: chat.SocketBufferSize, WriteBufferSize: chat.SocketBufferSize}

type RoomHandler struct {
	rooms map[string]*chat.Room
	*chat.Room
}

func (rh *RoomHandler) handleLeaveRoom(w http.ResponseWriter, req *http.Request) {
	//find client
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
	}
	data := objx.MustFromBase64(authCookie.Value)
	client := rh.Room.FindClient(data["userid"].(string))
	if client == nil {
		http.Error(w, errors.New("Cannot find client!").Error(), http.StatusUnauthorized)
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
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
	}
	data := objx.MustFromBase64(authCookie.Value)
	client := rh.Room.AddClient(socket, data)
	defer rh.Room.RemoveClient(client)
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
