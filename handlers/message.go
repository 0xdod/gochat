package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fibreactive/chat/chat"
	"github.com/fibreactive/chat/models"
)

type MessageHandler struct {
	models.MessageService
	models.RoomService
	models.UserService
}

func (mh *MessageHandler) Save(w http.ResponseWriter, r *http.Request) {
	data := CreateMessageJSON{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	m := &models.Message{
		Message: data.Message,
	}
	m.User = Get(r, "user").(*models.User)
	m.Room = Get(r, "room").(*models.Room)
	if err := mh.MessageService.Create(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// map every message to chat.message
// on load, make ajax request to get chats
// return chats and append html
func (mh *MessageHandler) List(w http.ResponseWriter, r *http.Request) {
	room := Get(r, "room").(*models.Room)
	session, _ := store.Get(r, "session.id")
	isPresent := session.Values["present"].(bool)
	messages := mh.RoomService.GetMessages(room)
	newMessages := make([]*chat.Message, len(messages))
	if !isPresent {
		w.WriteHeader(200)
		return
	}
	for i, message := range messages {
		message.User = mh.UserService.FindByID(uint(message.UserID))
		newMessage := chat.NewMessage(nil, message.User, message.Message)
		newMessages[i] = newMessage
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(newMessages)
	if err != nil {
		w.WriteHeader(422)
		return
	}
}
