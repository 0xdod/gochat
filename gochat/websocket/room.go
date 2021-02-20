package ws

// Room manages the routing of messages and clients
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
}

// Create a new chat room

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
	for {
		select {}
	}
}
