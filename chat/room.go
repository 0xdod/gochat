package chat

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

//type room models a single chat room
type room struct {
	//used to id a room
	name string
	//forward is a channel that holds incoming messages that should be forwarded to other clients
	forward chan *message
	//broadcast is a channel that holds incoming messages that should be forwarded to other clients except the sender
	broadcast chan *message
	// join is a channel for clients wishing to join the room.
	join chan *client
	// leave is a channel for clients wishing to leave the room.
	leave chan *client
	// clients holds all current clients in this room.
	clients map[*client]bool
	// tracer will receive trace information of activity in the room.
	// avatar is how avatar information will be obtained.
	//avatar Avatar
}

// newRoom makes a new room.
func NewRoom() *room {
	return &room{
		forward:   make(chan *message),
		broadcast: make(chan *message, 1),
		join:      make(chan *client),
		leave:     make(chan *client),
		clients:   make(map[*client]bool),
	}
}

func CreateRoom(name string) *room {
	return &room{
		name:      name,
		forward:   make(chan *message),
		broadcast: make(chan *message, 1),
		join:      make(chan *client),
		leave:     make(chan *client),
		clients:   make(map[*client]bool),
	}
}

//operation of the chat room
func (r *room) Run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
		case client := <-r.leave:
			// leaving
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// forward message to all clients
			for client := range r.clients {
				client.send <- msg

			}
		case msg := <-r.broadcast:
			// forward message to all clients except this one
			for client := range r.clients {
				if id, ok := client.userData["userid"]; ok {
					if id.(string) != msg.UserID {
						client.send <- msg
					}
				}

			}
		}
	}
}

//To get a web socket connection
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	newUserMsg := &message{
		Name:    "Admin",
		UserID:  client.userData["userid"].(string),
		Message: client.userData["name"].(string) + " just joined!",
		When:    time.Now(),
	}
	r.join <- client
	r.broadcast <- newUserMsg
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
