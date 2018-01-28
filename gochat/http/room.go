package http

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/0xdod/gochat/gochat"
	"github.com/0xdod/gochat/gochat/websocket"
)

func init() {
	go ws.GeneralRoom.Run()
}

var availableRooms = make(map[*ws.Room]bool)

func (s *Server) chat(w http.ResponseWriter, r *http.Request) {
	// does chat room exist?
	// if yes, does user already belong to room.
	// if yes add new participant
	// ?if no, retrieve old messages and render page (or do nothing)
	// if no room must be created.
	s.render(w, r, "chat.html", templateData{})
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
	room := ws.NewRoom()

	form := new(roomCreateForm)
	_ = parseForm(r, form)

	if err := form.validate(); err != nil {
		s.render(w, r, "create_room.html", templateData{
			"errors": err,
			"form":   form,
		})
		return
	}
	err := s.services.room.CreateRoom(context.Background(), form.create())
	if err != nil {
		s.serverError(w, err)
		return
	}
	availableRooms[room] = true
	addFlash(w, r, "success", fmt.Sprintf("Successfully created %q room", form.Name))
	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

func (s *Server) listRooms(w http.ResponseWriter, r *http.Request) {
	filter := gochat.RoomFilter{}
	rooms, _, err := s.services.room.FindRooms(context.Background(), filter)

	if err != nil {
		s.serverError(w, err)
		return
	}
	s.render(w, r, "room_list.html", templateData{
		"rooms": rooms,
	})
}
