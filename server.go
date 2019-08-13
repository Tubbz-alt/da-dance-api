package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/eveld/dda-api/protos/dda"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

// Server something
type Server struct {
	database *sqlx.DB
	router   *mux.Router
	logger   hclog.Logger
}

// NewServer creates a new server
func NewServer(logger hclog.Logger, router *mux.Router, database *sqlx.DB) *Server {
	server := &Server{
		logger:   logger,
		router:   router,
		database: database,
	}

	router.HandleFunc("/games", server.getGamesHandler).Methods(http.MethodGet)
	router.HandleFunc("/games/new", server.createGameHandler).Methods(http.MethodPost)
	router.HandleFunc("/games/{game}", server.getGameHandler).Methods(http.MethodGet)
	router.HandleFunc("/games/{game}/join", server.joinGameHandler).Methods(http.MethodPost)
	router.HandleFunc("/games/{game}/players/{player}/ready", server.readyGameHandler).Methods(http.MethodPost)
	router.HandleFunc("/games/{game}/start", server.startGameHandler).Methods(http.MethodPost)

	router.HandleFunc("/allocations", server.getAllocationsHandler).Methods(http.MethodGet)
	router.HandleFunc("/allocations/{allocation}/stop", server.stopAllocationHandler).Methods(http.MethodGet)

	return server
}

// func (s *Server) ClientStream(ctx context.Context, opts ...grpc.CallOption) (dda.DDAService_ClientStreamClient, error) {
// 	return dda.DDAService_ClientStreamClient{}, nil
// }

// GetGames fuck
func (s *Server) GetGames(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*dda.Games, error) {
	return nil, nil
}

// GetGame fuck
func (s *Server) GetGame(ctx context.Context, in *Game, opts ...grpc.CallOption) (*dda.Game, error) {
	return nil, nil
}

// JoinGame fuck
func (s *Server) JoinGame(ctx context.Context, in *Game, opts ...grpc.CallOption) (*dda.Game, error) {
	return nil, nil
}

// ReadyPlayer fuck
func (s *Server) ReadyPlayer(ctx context.Context, in *Game, opts ...grpc.CallOption) (*dda.Game, error) {
	return nil, nil
}

// StartGame fuck
func (s *Server) StartGame(ctx context.Context, in *Game, opts ...grpc.CallOption) (*Game, error) {
	return nil, nil
}

// GetAllocations fuck
func (s *Server) GetAllocations(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Allocations, error) {
	return nil, nil
}

// StopAllocation fuck
func (s *Server) StopAllocation(ctx context.Context, in *Allocation, opts ...grpc.CallOption) (*Allocation, error) {
	return nil, nil
}

// Start something
func (s *Server) Start() {
	s.logger.Info("Starting server", "address", *listenAddress)
	http.ListenAndServe(*listenAddress, s.router)
}

// GetGames gets all games
func (s *Server) GetGames() ([]Game, error) {
	var games []Game
	err := s.database.Select(&games, "SELECT * FROM games ORDER BY started DESC")
	if err == sql.ErrNoRows {
		s.logger.Error("No games found")
		return games, nil
	}

	if err != nil {
		s.logger.Error("Get games", "error", err)
		return games, err
	}

	return games, nil
}

// GetGame gets a specific game by ID
func (s *Server) GetGame(id string) (Game, error) {
	// Fetch game
	game := Game{
		ID: id,
	}

	gameQuery, err := s.database.PrepareNamed(
		`SELECT * FROM games
		WHERE id = :id`)
	if err != nil {
		s.logger.Error("Prepare game query", "error", err)
	}

	err = gameQuery.Get(&game, game)
	if err == sql.ErrNoRows {
		s.logger.Error("No game with that ID")
		return game, err
	}

	if err != nil {
		s.logger.Error("Query game", "error", err)
		return game, err
	}

	return game, nil
}

// Get all games
func (s *Server) getGamesHandler(w http.ResponseWriter, r *http.Request) {
	games, err := s.GetGames()
	if err != nil {
		//
	}

	out := json.NewEncoder(w)
	out.Encode(games)
}

// Create a new game
func (s *Server) createGameHandler(w http.ResponseWriter, r *http.Request) {
	gameID := uuid.New().String()
	playerID := uuid.New().String()

	game := Game{
		ID:        gameID,
		Song:      "first song",
		HomeID:    playerID,
		HomeScore: 0,
		HomeReady: false,
		AwayID:    "",
		AwayScore: 0,
		AwayReady: false,
		Started:   time.Time{},
		Finished:  time.Time{},
	}

	query, err := s.database.PrepareNamed(
		`INSERT INTO games (id, song, home_id, home_score, home_ready, away_id, away_score, away_ready, started, finished) 
		VALUES(:id, :song, :home_id, :home_score, :home_ready, :away_id, :away_score, :away_ready, :started, :finished) 
		RETURNING *`)
	if err != nil {
		s.logger.Error("Preparing game query", "error", err)
	}

	err = query.Get(&game, game)
	if err != nil {
		s.logger.Error("Querying game", "error", err)
	}

	out := json.NewEncoder(w)
	out.Encode(game)
}

// Join an existing game
func (s *Server) joinGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["game"]

	game, err := s.GetGame(gameID)
	if err != nil {
		//
	}

	if !game.Started.IsZero() {
		s.logger.Error("Game has started")
		return
	}

	if game.AwayID != "" {
		s.logger.Error("Game is full")
		return
	}

	game.AwayID = uuid.New().String()

	// Join game
	joinQuery, err := s.database.PrepareNamed(
		`UPDATE games
		SET away_id = :away_id
		WHERE id = :id
		RETURNING *`)
	if err != nil {
		s.logger.Error("Prepare game query", "error", err)
	}

	err = joinQuery.Get(&game, game)
	if err == sql.ErrNoRows {
		s.logger.Error("No game with that ID")
		return
	}

	if err != nil {
		s.logger.Error("Join game", "error", err)
		return
	}

	out := json.NewEncoder(w)
	out.Encode(game)
}

// Set player status to ready
func (s *Server) readyGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerID := vars["player"]

	gameID := vars["game"]

	game, err := s.GetGame(gameID)
	if err != nil {
		//
	}

	if playerID == game.HomeID {
		game.HomeReady = !game.HomeReady
	} else if playerID == game.AwayID {
		game.AwayReady = !game.AwayReady
	} else {
		s.logger.Error("Unknown player")
		return
	}

	// Set players ready
	readyQuery, err := s.database.PrepareNamed(
		`UPDATE games
		SET 
		home_ready = :home_ready,
		away_ready = :away_ready
		WHERE id = :id
		RETURNING *`)
	if err != nil {
		s.logger.Error("Prepare game query", "error", err)
	}

	err = readyQuery.Get(&game, game)
	if err == sql.ErrNoRows {
		s.logger.Error("No game with that ID")
		return
	}

	if err != nil {
		s.logger.Error("Set player ready", "error", err)
		return
	}

	out := json.NewEncoder(w)
	out.Encode(game)
}

// Get game status
func (s *Server) getGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	gameID := vars["game"]

	game, err := s.GetGame(gameID)
	if err != nil {
		//
	}

	out := json.NewEncoder(w)
	out.Encode(game)
}

// Start existing game
func (s *Server) startGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	gameID := vars["game"]

	game, err := s.GetGame(gameID)
	if err != nil {
		//
	}

	if !game.Started.IsZero() {
		s.logger.Error("Game has started")
		return
	}

	// Check if everyone is ready
	if game.HomeID != "" && !game.HomeReady {
		s.logger.Error("Host is not ready")
		return
	}

	if game.AwayID != "" && !game.AwayReady {
		s.logger.Error("Opponent is not ready")
		return
	}

	game.Started = time.Now()

	// Start the game
	startQuery, err := s.database.PrepareNamed(
		`UPDATE games
		SET 
		started = :started,
		WHERE id = :id
		RETURNING *`)
	if err != nil {
		s.logger.Error("Prepare game query", "error", err)
	}

	err = startQuery.Get(&game, game)
	if err == sql.ErrNoRows {
		s.logger.Error("No game with that ID")
		return
	}

	if err != nil {
		s.logger.Error("Start game", "error", err)
		return
	}

	out := json.NewEncoder(w)
	out.Encode(game)
}

// Get allocations
func (s *Server) getAllocationsHandler(w http.ResponseWriter, r *http.Request) {
	out := json.NewEncoder(w)
	out.Encode([]string{})
}

// Stop an allocation
func (s *Server) stopAllocationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	allocationID := vars["allocation"]

	out := json.NewEncoder(w)
	out.Encode(allocationID)
}
