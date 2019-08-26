package main

import (
	"fmt"

	"github.com/eveld/ddr-api/nomad"
	"github.com/eveld/ddr-api/server"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nicholasjackson/env"
)

var logger hclog.Logger

var listenAddress = env.String("LISTEN_ADDR", false, ":9090", "IP address and port to bind service to")
var databaseHost = env.String("POSTGRES_HOST", false, "localhost", "Host of the PostgreSQL server")
var databasePort = env.Int("POSTGRES_PORT", false, 5432, "Port of the PostgreSQL server")
var databaseUser = env.String("POSTGRES_USER", false, "", "Username of the PostgreSQL database")
var databasePassword = env.String("POSTGRES_PASSWORD", false, "", "Password of the PostgreSQL database")
var databaseName = env.String("POSTGRES_DATABASE", false, "", "Name of the PostgreSQL database")

func main() {
	logger = hclog.Default()

	env.Parse()

	router := mux.NewRouter()

	datasource := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", *databaseHost, *databasePort, *databaseUser, *databasePassword, *databaseName)
	database, err := sqlx.Connect("postgres", datasource)
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
