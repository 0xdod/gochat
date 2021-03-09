package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/0xdod/gochat/gochat"
	"github.com/0xdod/gochat/gochat/websocket"
	"github.com/gorilla/mux"
)

var availableRooms = make(map[int]*ws.Room)

func (s *Server) chat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inv := vars["invite"]
	// does chat room exist? : find by invite
	rf := gochat.RoomFilter{
		Invite: &inv,
	}
	rooms, _, err := s.FindRooms(context.TODO(), rf)
	if err != nil {
		s.serverError(w, err)
		return
	}
	if len(rooms) < 1 {
		addFlash(w, r, "error", "chat room does not exist")
		http.Redirect(w, r, "/rooms", http.StatusSeeOther)
		return
	}
	// if room exists, is user a participant? : find rooms where user is present.

	// if participant, enter chat, retrieve old messages
	// else ask to join [no message retrieval].
	// if room does not exist must be created.
	// retrieve all user rooms
	s.render(w, r, "chat.html", templateData{
		"room": rooms[0],
	})
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	socket, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("handleWS: ", err)
	}
	client := ws.NewClient(socket)
	// is room already running?
	//if yes retrieve room and add client
	// if no create new room and add client
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.clientError(w, http.StatusBadRequest)
		return
	}
	room, ok := availableRooms[id]
	if !ok {
		room = ws.NewRoom()
		availableRooms[id] = room
	}
	if !room.IsRunning() {
		go room.Run()
	}
	room.Join(client)
	go client.Write()
	client.Read()
	defer room.Leave(client)
}

func (s *Server) createRoom(w http.ResponseWriter, r *http.Request) {
	form := new(roomCreateForm)
	_ = parseForm(r, form)

	if err := form.validate(); err != nil {
		s.render(w, r, "create_room.html", templateData{
			"error": err,
			"form":  form,
		})
		return
	}
	user := UserFromContext(r.Context())
	room := form.create()
	room.CreatorID = user.ID
	room.AddParticipant(user).AddAdmin(user)
	link := room.InviteLink(r.Host)
	err := s.CreateRoom(context.TODO(), room)
	if err != nil {
		s.serverError(w, err)
		return
	}
	addFlash(w, r, "success", fmt.Sprintf("Successfully created %q room", form.Name))
	http.Redirect(w, r, link, http.StatusSeeOther)
}

func (s *Server) joinRoom(w http.ResponseWriter, r *http.Request) {
	form := new(joinRoomForm)
	_ = parseForm(r, form)

	if err := validateStruct(form); err != nil {
		s.clientError(w, http.StatusBadRequest)
		return
	}
	user := UserFromContext(r.Context())
	upd := gochat.RoomUpdate{
		Participants: []*gochat.User{user},
	}
	roomID, err := strconv.Atoi(form.ID)
	room, err := s.UpdateRoom(r.Context(), roomID, upd)
	if err != nil {
		s.serverError(w, err)
		return
	}
	addFlash(w, r, "success", fmt.Sprintf("Successfully joined %q room", room.Name))
	http.Redirect(w, r, room.Link, http.StatusSeeOther)
}

func (s *Server) listRooms(w http.ResponseWriter, r *http.Request) {
	// filter := gochat.RoomFilter{}
	// rooms, _, err := s.FindRooms(context.Background(), filter)

	// if err != nil {
	// 	s.serverError(w, err)
	// 	return
	// }
	s.render(w, r, "room_list.html", templateData{
		//"allRooms": rooms,
	})
}
