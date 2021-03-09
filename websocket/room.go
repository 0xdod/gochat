package ws

// Room manages the clients and routing of messages.
type Room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to other clients.
	forward chan *Message

	// join is a channel for clients wishing to join the room.
	join chan *Client

	// leave is a channel for clients wishing to leave the room.
	leave chan *Client

	// clients manages all clients in this room.
	clients map[*Client]bool

	status string
}

// NewRoom create a new chat room
func NewRoom() *Room {
	return &Room{
		forward: make(chan *Message),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		clients: make(map[*Client]bool),
	}
}

// Run the operation of the chat room
func (r *Room) Run() {
	r.status = "running"
	for {
		select {
		case client := <-r.join:
			client.room = r
			r.clients[client] = true
		case client := <-r.leave:
			close(client.send)
			delete(r.clients, client)
		case msg := <-r.forward:
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

func (r *Room) IsRunning() bool {
	return r.status == "running"
}

// Join a chat room
func (r *Room) Join(client *Client) {
	r.join <- client
}

// Leave a chat room
func (r *Room) Leave(client *Client) {
	r.leave <- client
}

// GeneralRoom is a general room.
var GeneralRoom = &Room{
	forward: make(chan *Message),
	join:    make(chan *Client),
	leave:   make(chan *Client),
	clients: make(map[*Client]bool),
}
