package handlers

import (
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

func (rh *RoomHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
