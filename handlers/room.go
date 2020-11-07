package handlers

// import (
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/gorilla/websocket"
// 	"github.com/stretchr/objx"
// )

// const (
// 	socketBufferSize  = 1024
// 	messageBufferSize = 256
// )

// //To get a web socket connection
// var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

// type roomHandler struct {
// 	room *chat.Room
// }

// func (rh *roomHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	socket, err := upgrader.Upgrade(w, req, nil)
// 	if err != nil {
// 		log.Fatal("ServeHTTP:", err)
// 		return
// 	}
// 	authCookie, err := req.Cookie("auth")
// 	if err != nil {
// 		log.Fatal("Failed to get auth cookie:", err)
// 	}
// 	client := &chat.Client{
// 		Socket:   socket,
// 		Send:     make(chan *chat.Message, messageBufferSize),
// 		Room:     rh.room,
// 		UserData: objx.MustFromBase64(authCookie.Value),
// 	}
// 	newUserMsg := &chat.Message{
// 		Name:    "Admin",
// 		Message: client.userData["name"].(string) + " just joined!",
// 		When:    time.Now(),
// 	}
// 	rh.room.join <- client
// 	rh.room.forward <- newUserMsg
// 	defer func() { rh.room.leave <- client }()
// 	go client.write()
// 	client.read()
// }
