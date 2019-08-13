package main

import (
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nicholasjackson/env"
)

var logger hclog.Logger

var listenAddress = env.String("LISTEN_ADDR", false, ":9090", "IP address and port to bind service to")

// Game represents a game of DDA.
type Game struct {
	ID        string    `json:"id" db:"id"`
	Song      string    `json:"song" db:"song"`
	HomeID    string    `json:"home_id" db:"home_id"`
	HomeScore int       `json:"home_score" db:"home_score"`
	HomeReady bool      `json:"home_ready" db:"home_ready"`
	AwayID    string    `json:"away_id" db:"away_id"`
	AwayScore int       `json:"away_score" db:"away_score"`
	AwayReady bool      `json:"away_ready" db:"away_ready"`
	Started   time.Time `json:"started" db:"started"`
	Finished  time.Time `json:"finished" db:"finished"`
}

func main() {
	logger = hclog.Default()

	env.Parse()

	router := mux.NewRouter()
	database, err := sqlx.Connect("postgres", "user=secret_user password=secret_password dbname=dda sslmode=disable")
	if err != nil {
		logger.Error("Connecting to postgres", err)
	}

	server := NewServer(logger, router, database)
	server.Start()
}
