package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/fibreactive/chat/models"
	"github.com/jinzhu/gorm"
)

var mh *MessageHandler
var modifiedTestMessages []*models.Message
var testMessages = []*models.Message{
	&models.Message{
		Model: gorm.Model{
			ID: 1,
		},
		UserID:  1,
		RoomID:  1,
		Message: "Hi",
	},
	&models.Message{
		Model: gorm.Model{
			ID: 2,
		},
		UserID:  2,
		RoomID:  1,
		Message: "Hi",
	},
	&models.Message{
		Model: gorm.Model{
			ID: 3,
		},
		UserID:  1,
		RoomID:  2,
		Message: "Hi",
	},
	&models.Message{
		Model: gorm.Model{
			ID: 4,
		},
		UserID:  3,
		RoomID:  2,
		Message: "Hi",
	},
}

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	// prepare db and seed data
	tms := &TestMessageService{}
	trs := &TestRoomService{}
	tus := &TestUserService{}
	mh = &MessageHandler{tms, trs, tus}

}

func tearDown() {
	// destroy db data
}

func TestCreate(t *testing.T) {
	postData := strings.NewReader(`{"message":"Beast of no nation", "roomID": 1, "userID":2}`)
	request, _ := http.NewRequest("POST", "/messages", postData)
	writer := httptest.NewRecorder()
	request = Set(request, "room", testRooms[0])
	request = Set(request, "user", testUsers[0])
	mh.Create(writer, request)
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	if len(modifiedTestMessages)-len(testMessages) != 1 {
		t.Errorf("Expected new data set of %v, but got %v", len(testMessages), len(modifiedTestMessages))
	}
}

func TestListForOldMember(t *testing.T) {
	// previous messages are retrieved for old members
	request, _ := http.NewRequest("GET", "/messages", nil)
	writer := httptest.NewRecorder()
	session, _ := store.Get(request, "session.id")
	session.Values["present"] = true
	if err := session.Save(request, writer); err != nil {
		t.Error(err.Error())
	}
	request = Set(request, "room", testRooms[0])
	mh.List(writer, request)
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	if ct := writer.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content type is %v", ct)
	}
	var msgs []*models.Message
	json.Unmarshal(writer.Body.Bytes(), &msgs)
	if len(msgs) != 2 {
		t.Errorf("Expected 2 messages, but got %v", len(msgs))
	}
}

func TestListForNewMember(t *testing.T) {
	// no previous messages are retrieved for a new member
	request, _ := http.NewRequest("GET", "/messages", nil)
	writer := httptest.NewRecorder()
	session, _ := store.Get(request, "session.id")
	session.Values["present"] = false
	if err := session.Save(request, writer); err != nil {
		t.Error(err.Error())
	}
	request = Set(request, "room", testRooms[0])
	mh.List(writer, request)
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	if ct := writer.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content type is %v", ct)
	}
	var ms []*models.Message
	json.Unmarshal(writer.Body.Bytes(), &ms)
	if len(ms) != 0 {
		t.Errorf("Expected 0 messages, but got %v", len(ms))
	}
}

type TestMessageService struct{}

func (tms *TestMessageService) GetAll() []*models.Message {
	return nil
}

func (tms *TestMessageService) GetByRoom(r *models.Room) []*models.Message {
	roomMessages := make([]*models.Message, 0)
	for _, tm := range testMessages {
		if r.Model.ID == tm.RoomID {
			roomMessages = append(roomMessages, tm)
		}
	}
	return roomMessages
}

func (tms *TestMessageService) GetByUser(r *models.User) []*models.Message {
	return nil
}

func (tms *TestMessageService) FindOne(m models.Map) *models.Message {
	return nil
}

func (tms *TestMessageService) FindMany(m models.Map) []*models.Message {
	return nil
}

func (tms *TestMessageService) Create(m *models.Message) error {
	newLength := len(testMessages) + 1
	modifiedTestMessages = make([]*models.Message, newLength)
	copy(modifiedTestMessages, testMessages)
	modifiedTestMessages[newLength-1] = m
	return nil
}

func (tms *TestMessageService) Update(m *models.Message) error {
	return nil
}

func (tms *TestMessageService) Delete(id uint) error {
	return nil
}
