package http

import (
	"log"
	"net/http"

	"github.com/0xdod/gochat/gochat/http/templates"
	"github.com/0xdod/gochat/gochat/websocket"
)

func init() {
	go ws.GeneralRoom.Run()
}

func index(w http.ResponseWriter, r *http.Request) {
	templates.Render(w, "chat.html", r)
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	socket, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("handleWS: ", err)
	}
	client := ws.NewClient(socket)
	ws.GeneralRoom.Join(client)
	go client.Write()
	client.Read()
}
