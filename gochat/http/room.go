package http

import (
	"log"
	"net/http"

	"github.com/0xdod/gochat/gochat/websocket"
)

func init() {
	go ws.GeneralRoom.Run()
}

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
