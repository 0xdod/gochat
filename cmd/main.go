package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/0xdod/gochat/http"
	"github.com/0xdod/gochat/store"
)

// Move these to env and conf
const (
	host     = "localhost"
	port     = 5432
	user     = "damilola"
	password = "Omonefe97"
	dbname   = "gochat"
)

type Config struct {
	Addr string
	DSN  string
}

func getConfig() *Config {
	cfg := new(Config)
	localDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	flag.StringVar(&cfg.Addr, "addr", getEnvWithDefault("PORT", "4000"), "HTTP server address")
	flag.StringVar(&cfg.DSN, "dsn", getEnvWithDefault("DATABASE_URL", localDSN), "Data source name")
	flag.Parse()
	return cfg
}

func main() {
	cfg := getConfig()
	db, err := store.ConnectToDB(cfg.DSN)
	if err != nil {
		panic(err)
	}
	s := http.NewServer()
	s.LoadTemplates()
	s.CreateServices(db)
	s.InfoLog.Printf("Server running on port %s", cfg.Addr)
	s.Run(cfg.Addr)
}

func getEnvWithDefault(name, def string) string {
	p := os.Getenv(name)
	if p == "" {
		p = def
	}
	if name == "PORT" {
		return ":" + p
	}
	return p
}
