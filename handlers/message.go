package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fibreactive/chat/chat"
	"github.com/fibreactive/chat/models"
	"github.com/stretchr/objx"
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
	messages := mh.RoomService.GetMessages(room)
	newMessages := make([]*chat.Message, len(messages))
	for _, message := range messages {
		message.User = mh.UserService.FindByID(uint(message.UserID))
		newMessage := chat.NewMessageAlt(nil, message.User, message.Message)
		newMessages = append(newMessages, newMessage)
	}
	data := objx.MSI()
	data.Set("user", Get(r, "user").(*models.User))
	data.Set("room", room)
	data.Set("messages", newMessages)
	renderOne(w, "message_list_ajax", data)
}
