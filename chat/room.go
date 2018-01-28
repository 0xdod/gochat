package chat

import (
	"time"

	"github.com/gorilla/websocket"
)

//type room models a single chat room
type Room struct {
	//forward is a channel that holds incoming messages that should be forwarded to other clients
	forward chan *Message
	//broadcast is a channel that holds incoming messages that should be forwarded to other clients except the sender
	// join is a channel for clients wishing to join the room.
	join chan *Client
	// leave is a channel for clients wishing to leave the room.
	leave chan *Client
	// clients holds all current clients in this room.
	clients map[*Client]bool
	// track number of clients present
	Nclients int
	// how to get extra room info
	DataFinder
}

// we need room name, id

func NewRoom(df DataFinder) *Room {
	return &Room{
		forward:    make(chan *Message),
		join:       make(chan *Client),
		leave:      make(chan *Client),
		clients:    make(map[*Client]bool),
		DataFinder: df,
	}
}

//operation of the chat room
func (r *Room) Run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
			r.Nclients++
			if !client.IsPresent {
				msg := generateAdminMessage(client, NEW_USER)
				r.sendMessageToClient(client, msg)
				msg = generateAdminMessage(client, USER_JOINED)
				r.sendMessageToClientsExcept(client, msg)
			}
			client.IsPresent = true
		case client := <-r.leave:
			// leaving
			r.Nclients--
			delete(r.clients, client)
			close(client.send)
			r.sendMessageToClientsExcept(client, generateAdminMessage(client, USER_LEFT))
		case msg := <-r.forward:
			// forward message to all clients
			for client := range r.clients {
				client.send <- msg
			}
		}
		// if no more client present stop the room.
		if r.Nclients < 1 {
			return
		}
	}
}

func (r *Room) CountClients() int {
	var clients []*Client
	for c := range r.clients {
		clients = append(clients, c)
	}
	r.Nclients = len(clients)
	return r.Nclients
}

func (r *Room) AddClient(client *Client) {
	// find client, if present remove
	if c := r.FindClient(client.GetIntID()); c != nil {
		r.RemoveClient(c)
	}
	client.room = r
	r.join <- client
}

func (r *Room) RemoveClient(client *Client) {
	closeMessage := websocket.FormatCloseMessage(websocket.CloseGoingAway, "client has left")
	client.socket.WriteControl(websocket.CloseMessage, closeMessage, time.Now().Add(time.Millisecond*15))
	if _, ok := r.clients[client]; ok {
		r.leave <- client
	}
}

func (r *Room) FindClient(id int) *Client {
	for client, _ := range r.clients {
		if client.GetIntID() == id {
			return client
		}
	}
	return nil
}

func (r *Room) sendMessageToClient(c *Client, msg *Message) {
	c.send <- msg
}

func (r *Room) sendMessageToClientsExcept(ignore *Client, msg *Message) {
	for client := range r.clients {
		if client != ignore {
			client.send <- msg
		}
	}
}
