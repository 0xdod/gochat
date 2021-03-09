package main

import (
	"log"
	"net/http"
	"os"

	"github.com/urfave/negroni"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	n := negroni.Classic()
	MapRoutes(n)

	log.Println("Starting server on port 9000")
	http.ListenAndServe(":"+port, n)
}
