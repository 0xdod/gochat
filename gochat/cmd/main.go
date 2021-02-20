package main

import (
	"log"

	"github.com/0xdod/gochat/gochat/http"
)

func main() {

	s := http.NewServer()

	log.Println("Server running on port :8080")
	s.Run(":8080")
}
