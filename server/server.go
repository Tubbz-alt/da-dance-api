package server

import (
	"net/http"

	"github.com/eveld/ddr-api/nomad"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
)

// Server something
type Server struct {
	database *sqlx.DB
	router   *mux.Router
	logger   hclog.Logger
	nomad    *nomad.Nomad
}

// NewServer creates a new server
func NewServer(logger hclog.Logger, router *mux.Router, database *sqlx.DB, nomad *nomad.Nomad) *Server {
	server := &Server{
		logger:   logger,
		router:   router,
		database: database,
		nomad:    nomad,
	}

	router.HandleFunc("/games", server.getGamesHandler).Methods(http.MethodGet)
	router.HandleFunc("/games/new", server.createGameHandler).Methods(http.MethodPost)
	router.HandleFunc("/games/{game}", server.getGameHandler).Methods(http.MethodGet)
	router.HandleFunc("/games/{game}", server.deleteGameHandler).Methods(http.MethodDelete)
	router.HandleFunc("/games/{game}/join", server.joinGameHandler).Methods(http.MethodPost)
	router.HandleFunc("/games/{game}/leave", server.leaveGameHandler).Queries("player", "{player}").Methods(http.MethodPost)
	router.HandleFunc("/games/{game}/ready", server.readyGameHandler).Queries("player", "{player}").Methods(http.MethodPost)
	router.HandleFunc("/games/{game}/start", server.startGameHandler).Methods(http.MethodPost)

	router.HandleFunc("/allocations", server.getAllocationsHandler).Queries("player", "{player}", "count", "{count}").Methods(http.MethodGet)
	router.HandleFunc("/allocations/{allocation}/stop", server.stopAllocationHandler).Methods(http.MethodPost)

	return server
}

// Start something
func (s *Server) Start(address string) {
	s.logger.Info("Starting server", "address", address)
	http.ListenAndServe(address, s.router)
}
