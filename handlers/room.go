package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/fibreactive/chat/models"
	"github.com/stretchr/objx"

	"github.com/fibreactive/chat/chat"

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
	rooms := rh.RoomService.GetAll()
	if len(rooms) == 0 {
		rh.Rooms[0] = chat.NewRoom(nil)
		return
	}
	for _, room := range rooms {
		if _, exists := rh.Rooms[int(room.ID)]; !exists {
			rh.Rooms[int(room.ID)] = chat.NewRoom(room)
		}
	}

	for _, r := range rh.Rooms {
		if r.CountClients() > 0 {
			go r.Run()
		}
	}
}

func (rh *RoomHandler) ListOrCreate(w http.ResponseWriter, req *http.Request) {
	user, ok := Get(req, "user").(*models.UserModel)
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
		roomDB := &models.RoomModel{
			Name:        form.Name,
			Description: form.Description,
		}
		roomDB.Admins = append(roomDB.Admins, user)
		rh.RoomService.Create(roomDB)
		room := chat.NewRoom(roomDB)
		rh.Rooms[int(roomDB.ID)] = room
		http.Redirect(w, req, "/chat/"+roomDB.Link, http.StatusSeeOther)
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
	vars := mux.Vars(req)
	rID := vars["id"]
	id, _ := strconv.Atoi(rID)
	//find client
	user, ok := Get(req, "user").(*models.UserModel)
	if !ok {
		http.Error(w, errors.New("Cannot find client!").Error(), http.StatusForbidden)
		return
	}
	room := rh.Rooms[id]
	roomDB := rh.RoomService.FindByID(uint(id))
	rh.RoomService.RemoveParticipant(roomDB, user)
	client := room.FindClient(user.ID)
	if client == nil {
		http.Error(w, errors.New("Cannot find client!").Error(), http.StatusForbidden)
		return
	}
	room.RemoveClient(client)
	w.WriteHeader(http.StatusOK)
	return
}

func (rh *RoomHandler) Chat(w http.ResponseWriter, req *http.Request) {
	user, ok := Get(req, "user").(*models.UserModel)
	if !ok {
		http.Error(w, errors.New("Cannot find client!").Error(), http.StatusForbidden)
		return
	}
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	vars := mux.Vars(req)
	roomID := vars["id"]
	id, _ := strconv.Atoi(roomID)
	room := rh.Rooms[id]
	if room.CountClients() < 1 {
		go room.Run()
	}
	client := chat.NewClient(socket, user)
	defer room.RemoveClient(client)
	room.AddClient(client)
	go client.Write()
	client.Read()
}

//func (rh *RoomHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {}
