package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/eveld/ddr-api/models"
	"github.com/eveld/ddr-api/nomad"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	router.HandleFunc("/games/{game}/join", server.joinGameHandler).Methods(http.MethodPost)
	router.HandleFunc("/games/{game}/players/{player}/ready", server.readyGameHandler).Methods(http.MethodPost)
	router.HandleFunc("/games/{game}/start", server.startGameHandler).Methods(http.MethodPost)

	router.HandleFunc("/allocations", server.getAllocationsHandler).Methods(http.MethodGet)
	router.HandleFunc("/allocations/{allocation}/stop", server.stopAllocationHandler).Methods(http.MethodGet)

	return server
}

// Start something
func (s *Server) Start(address string) {
	s.logger.Info("Starting server", "address", address)
	http.ListenAndServe(address, s.router)
}

// GetGames gets all games
func (s *Server) GetGames() ([]models.Game, error) {
	var games []models.Game
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
func (s *Server) GetGame(id string) (models.Game, error) {
	// Fetch game
	game := models.Game{
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

// CreateGame creates a new game
func (s *Server) CreateGame(game models.Game) (models.Game, error) {
	query, err := s.database.PrepareNamed(
		`INSERT INTO games (id, song, home_id, home_score, home_ready, away_id, away_score, away_ready, started, finished)
		VALUES(:id, :song, :home_id, :home_score, :home_ready, :away_id, :away_score, :away_ready, :started, :finished)
		RETURNING *`)
	if err != nil {
		return game, err
	}

	err = query.Get(&game, game)
	if err != nil {
		return game, err
	}

	return game, nil
}

// UpdateGame updates an existing game
func (s *Server) UpdateGame(game models.Game) (models.Game, error) {
	// Join game
	joinQuery, err := s.database.PrepareNamed(
		`UPDATE games
		SET
			song = :song, 
			home_id = :home_id, 
			home_score = :home_score, 
			home_ready = :home_ready, 
			away_id = :away_id, 
			away_score = :away_score, 
			away_ready = :away_ready, 
			started = :started, 
			finished = :finished
		WHERE id = :id
		RETURNING *`)
	if err != nil {
		return game, err
	}

	err = joinQuery.Get(&game, game)
	if err != nil {
		return game, err
	}

	return game, nil
}

// Get all games
func (s *Server) getGamesHandler(w http.ResponseWriter, r *http.Request) {
	games, err := s.GetGames()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := json.NewEncoder(w)
	out.Encode(games)
}

// Create a new game
func (s *Server) createGameHandler(w http.ResponseWriter, r *http.Request) {
	gameID := uuid.New().String()
	playerID := uuid.New().String()

	game := models.Game{
		ID:        gameID,
		Song:      "first song",
		HomeID:    playerID,
		HomeScore: 0,
		HomeReady: false,
		AwayID:    "",
		AwayScore: 0,
		AwayReady: false,
		Started:   0,
		Finished:  0,
	}

	game, err := s.CreateGame(game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if game.Started != 0 {
		http.Error(w, "Game has started", http.StatusInternalServerError)
		return
	}

	if game.AwayID != "" {
		http.Error(w, "Game is full", http.StatusInternalServerError)
		return
	}

	game.AwayID = uuid.New().String()
	game, err = s.UpdateGame(game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if playerID == game.HomeID {
		game.HomeReady = !game.HomeReady
	} else if playerID == game.AwayID {
		game.AwayReady = !game.AwayReady
	} else {
		http.Error(w, "Unknown player", http.StatusInternalServerError)
		return
	}

	game, err = s.UpdateGame(game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if game.Started != 0 {
		http.Error(w, "Game has started", http.StatusInternalServerError)
		return
	}

	// Check if everyone is ready
	if game.HomeID != "" && !game.HomeReady {
		http.Error(w, "Host is not ready", http.StatusInternalServerError)
		return
	}

	if game.AwayID != "" && !game.AwayReady {
		http.Error(w, "Opponent is not ready", http.StatusInternalServerError)
		return
	}

	game.Started = time.Now().Unix()
	game, err = s.UpdateGame(game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
