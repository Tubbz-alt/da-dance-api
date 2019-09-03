package main

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/hashicorp/da-dance-api/nomad"
	"github.com/hashicorp/da-dance-api/server"
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nicholasjackson/env"
)

var logger hclog.Logger

var listenAddress = env.String("LISTEN_ADDR", false, ":9090", "IP address and port to bind service to")
var databaseHost = env.String("POSTGRES_HOST", false, "localhost", "Host of the PostgreSQL server")
var databasePort = env.Int("POSTGRES_PORT", false, 5432, "Port of the PostgreSQL server")
var databaseUser = env.String("POSTGRES_USER", true, "", "Username of the PostgreSQL database")
var databasePassword = env.String("POSTGRES_PASSWORD", true, "", "Password of the PostgreSQL database")
var databaseName = env.String("POSTGRES_DATABASE", true, "", "Name of the PostgreSQL database")
var nomadJobID = env.String("NOMAD_JOB_ID", true, "", "ID of a Nomad job for whacking")

func main() {
	logger = hclog.Default()

	env.Parse()

	router := mux.NewRouter()

	datasource := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", *databaseHost, *databasePort, *databaseUser, *databasePassword, *databaseName)
	database, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		logger.Error("Connecting to postgres", "error", err)
	}

	nomad, err := nomad.Connect(*nomadJobID)
	if err != nil {
		logger.Error("Connecting to nomad", "error", err)
	}

	s := server.NewServer(logger, router, database, nomad)
	s.Start(*listenAddress)
}
