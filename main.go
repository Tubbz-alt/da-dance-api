package main

import (
	"github.com/eveld/ddr-api/server"
	"github.com/eveld/ddr-api/nomad"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nicholasjackson/env"
)

var logger hclog.Logger

var listenAddress = env.String("LISTEN_ADDR", false, ":9090", "IP address and port to bind service to")

func main() {
	logger = hclog.Default()

	env.Parse()

	router := mux.NewRouter()
	database, err := sqlx.Connect("postgres", "user=secret_user password=secret_password dbname=dda sslmode=disable")
	if err != nil {
		logger.Error("Connecting to postgres", "error", err)
	}

	nomad, err := nomad.Connect()
	if err != nil {
		logger.Error("Connecting to nomad", "error", err)
	}

	s := server.NewServer(logger, router, database, nomad)
	s.Start(*listenAddress)
}
