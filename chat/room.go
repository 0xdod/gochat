package chat

import "github.com/fibreactive/chat/models"

//type room models a single chat room
type Room struct {
	//forward is a channel that holds incoming messages that should be forwarded to other clients
	forward chan *message
	//broadcast is a channel that holds incoming messages that should be forwarded to other clients except the sender
	// join is a channel for clients wishing to join the room.
	join chan *Client
	// leave is a channel for clients wishing to leave the room.
	leave chan *Client
	// clients holds all current clients in this room.
	clients map[*Client]bool
	// tracer will receive trace information of activity in the room.
	// avatar is how avatar information will be obtained.
	//avatar Avatar
	room *models.RoomModel
}

func NewRoom(room *models.RoomModel) *Room {
	return &Room{
		forward: make(chan *message),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		clients: make(map[*Client]bool),
		room:    room,
	}
}

//operation of the chat room
func (r *Room) Run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
			msg := generateAdminMessage(client, NEW_USER)
			r.sendMessageToClient(client, msg)
			msg = generateAdminMessage(client, USER_JOINED)
			r.sendMessageToClientsExcept(client, msg)
		case client := <-r.leave:
			// leaving
			delete(r.clients, client)
			close(client.send)
			r.sendMessageToClientsExcept(client, generateAdminMessage(client, USER_LEFT))
		case msg := <-r.forward:
			// forward message to all clients
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

func (r *Room) GetRoomModel() *models.RoomModel {
	if r.room == nil {
		return nil
	}
	return r.room
}

func (r *Room) CountClients() int {
	var clients []*Client
	for c := range r.clients {
		clients = append(clients, c)
	}
	return len(clients)
}

func (r *Room) AddClient(client *Client) {
	client.room = r
	for c, v := range r.clients {
		if c.user.ID == client.user.ID && v {
			r.RemoveClient(c)
		}
	}
	if _, ok := r.clients[client]; !ok {
		r.join <- client
	}
}

func (r *Room) RemoveClient(client *Client) {
	if _, ok := r.clients[client]; ok {
		r.leave <- client
	}
}

func (r *Room) FindClient(id uint) *Client {
	for client := range r.clients {
		if client.user.ID == id {
			return client
		}
	}
	return nil
}

func (r *Room) sendMessageToClient(c *Client, msg *message) {
	c.send <- msg
}

func (r *Room) sendMessageToClientsExcept(ignore *Client, msg *message) {
	for client := range r.clients {
		if client != ignore {
			client.send <- msg
		}
	}
}
