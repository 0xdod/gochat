package http

import (
	"log"
	"net/http"

	"github.com/0xdod/gochat/gochat/websocket"
)

func init() {
	go ws.GeneralRoom.Run()
}

var availableRooms map[*ws.Room]bool

func (s *Server) chat(w http.ResponseWriter, r *http.Request) {
	s.render(w, "chat.html", r)
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	socket, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("handleWS: ", err)
	}
	client := ws.NewClient(socket)
	ws.GeneralRoom.Join(client)
	go client.Write()
	client.Read()
}

func (s *Server) createRoom(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		room := ws.NewRoom()
		availableRooms[room] = true
		go room.Run()
		http.Redirect(w, r, "/chat", 301)
		return
	}
	s.render(w, "room_list.html", r)

}

func (s *Server) roomList(w http.ResponseWriter, r *http.Request) {
	s.render(w, "room_list.html", r)
}
