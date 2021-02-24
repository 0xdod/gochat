package main

import (
	"fmt"
	"log"

	"github.com/0xdod/gochat/gochat/gorm"
	"github.com/0xdod/gochat/gochat/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "damilola"
	password = "Omonefe97"
	dbname   = "gochat"
)

func main() {
	localDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.ConnectToDB(localDSN)
	if err != nil {
		panic(err)
	}
	s := http.NewServer()
	s.LoadTemplates()
	s.CreateServices(db)
	log.Println("Server running on port :8080")
	s.Run(":8080")
}
