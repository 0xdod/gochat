package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/stretchr/objx"

	"github.com/fibreactive/golang-rtc/chat"
	"github.com/fibreactive/golang-rtc/models"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

//To get a web socket connection
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  chat.SocketBufferSize,
	WriteBufferSize: chat.SocketBufferSize,
}

type RoomHandler struct {
	models.RoomService
	models.UserService
	Rooms map[int]*chat.Room
}

func (rh *RoomHandler) PopulateRooms() {
	rh.Rooms = make(map[int]*chat.Room)
	roomSlice := rh.RoomService.GetAll()
	if len(roomSlice) == 0 {
		return
	}
	for _, room := range roomSlice {
		rh.Rooms[int(room.ID)] = chat.NewRoom(room)
	}
}

func (rh *RoomHandler) getChatRoom(id int) *chat.Room {
	if room, exists := rh.Rooms[id]; exists {
		return room
	}
	return nil
}

func (rh *RoomHandler) ListOrCreate(w http.ResponseWriter, req *http.Request) {
	user, ok := Get(req, "user").(*models.User)
	if !ok {
		http.Error(w, errors.New("Cannot find client!").Error(), http.StatusForbidden)
		return
	}
	if req.Method == "POST" {
		//redirect to chat page
		form := CreateRoomForm{}
		if err := parseForm(req, &form); err != nil {
			panic(err)
		}
		room := &models.Room{
			Name:        form.Name,
			Description: form.Description,
		}
		room.Admins = append(room.Admins, user)
		rh.RoomService.Create(room)
		rh.Rooms[int(room.ID)] = chat.NewRoom(room)
		http.Redirect(w, req, "/chat/"+room.Link, http.StatusSeeOther)
		return
	}
	m := objx.MSI()
	m.Set("rooms", rh.RoomService.GetAll())
	m.Set("MyRooms", rh.UserService.GetRooms(user))
	render(w, "room_list.html", m)
}

func (rh *RoomHandler) Edit(w http.ResponseWriter, req *http.Request) {

}

func (rh *RoomHandler) Leave(w http.ResponseWriter, req *http.Request) {
	user, ok := Get(req, "user").(*models.User)
	if !ok {
		http.Error(w, errors.New("Cannot find user!").Error(), http.StatusForbidden)
		return
	}
	room, ok := Get(req, "room").(*models.Room)
	if !ok {
		http.Error(w, errors.New("Cannot find room!").Error(), http.StatusForbidden)
		return
	}
	vars := mux.Vars(req)
	roomId, _ := strconv.Atoi(vars["id"])
	chatRoom := rh.getChatRoom(roomId)
	if chatRoom == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	client := chatRoom.FindClient(int(user.ID))
	if client == nil {
		http.Error(w, errors.New("Cannot find client!").Error(), http.StatusForbidden)
		return
	}
	chatRoom.RemoveClient(client)
	rh.RoomService.RemoveParticipant(room, user)
	w.WriteHeader(http.StatusOK)
	return
}

func (rh *RoomHandler) ServeWs(w http.ResponseWriter, req *http.Request) {
	user, ok := Get(req, "user").(*models.User)
	if !ok {
		http.Error(w, errors.New("Cannot find client!").Error(), http.StatusForbidden)
		return
	}
	session, _ := store.Get(req, "session.id")
	isPresent := session.Values["present"].(bool)
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := chat.NewClient(socket, user, isPresent)
	vars := mux.Vars(req)
	rmID := vars["id"]
	id, _ := strconv.Atoi(rmID)
	room := rh.getChatRoom(id)
	if room == nil {
		return
	}
	if room.Nclients < 1 {
		go room.Run()
	}
	room.AddClient(client)
	go client.Write()
	go client.Read()
}

//func (rh *RoomHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {}
