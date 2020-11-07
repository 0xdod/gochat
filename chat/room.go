package chat

import (
	"time"

	"github.com/gorilla/websocket"
)

//type room models a single chat room
type Room struct {
	//used to id a room
	name string
	//forward is a channel that holds incoming messages that should be forwarded to other clients
	forward chan *message
	//broadcast is a channel that holds incoming messages that should be forwarded to other clients except the sender
	broadcast chan *message
	// join is a channel for clients wishing to join the room.
	join chan *Client
	// leave is a channel for clients wishing to leave the room.
	leave chan *Client
	// clients holds all current clients in this room.
	clients map[*Client]bool
	// tracer will receive trace information of activity in the room.
	// avatar is how avatar information will be obtained.
	//avatar Avatar
}

func CreateRoom(name string) *Room {
	return &Room{
		name:      name,
		forward:   make(chan *message),
		broadcast: make(chan *message, 1),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		clients:   make(map[*Client]bool),
	}
}

//operation of the chat room
func (r *Room) Run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
			msgBody := client.userData["name"].(string) + ", welcome to " + client.room.name + " chat room"
			id := client.userData["userid"].(string)
			client.send <- adminWelcomeMessage(msgBody)
			msgBody = client.userData["name"].(string) + " has joined!"
			r.broadcast <- adminBroadcastMessage(id, msgBody)
		case client := <-r.leave:
			// leaving
			close(client.send)
			delete(r.clients, client)
			msgBody := client.userData["name"].(string) + " has left"
			id := client.userData["userid"].(string)
			r.broadcast <- adminBroadcastMessage(id, msgBody)
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

func (r *Room) AddClient(s *websocket.Conn, data map[string]interface{}) *Client {
	client := &Client{
		socket:   s,
		room:     r,
		userData: data,
		send:     make(chan *message, MessageBufferSize),
	}
	if _, ok := r.clients[client]; !ok {
		r.join <- client
	}
	return client
}

func (r *Room) FindClient(id string) *Client {
	for client := range r.clients {
		if userID := client.userData["userid"]; userID.(string) == id {
			return client
		}
	}
	return nil
}

func (r *Room) RemoveClient(client *Client) {
	if _, ok := r.clients[client]; ok {
		r.leave <- client
	}
}

func adminBroadcastMessage(userId, msg string) *message {
	return &message{
		UserID:  userId,
		Name:    "Admin",
		Message: msg,
		When:    time.Now(),
	}
}

func adminWelcomeMessage(msg string) *message {
	return &message{
		UserID:  "Admin",
		Name:    "Admin",
		Message: msg,
		When:    time.Now(),
	}
}
