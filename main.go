package main

import (
	"log"
	"net/http"

	"github.com/urfave/negroni"
)

func main() {
	n := negroni.Classic()
	MapRoutes(n)
	log.Println("Starting server on port 9000")
	http.ListenAndServe(":9000", n)
}
