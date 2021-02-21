package http

import (
	"net/http"

	"github.com/0xdod/gochat/gochat/http/templates"
)

func index(w http.ResponseWriter, r *http.Request) {
	templates.Render(w, "chat.html", r)
}
