package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eveld/ddr-api/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

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

// Delete an existing game
func (s *Server) deleteGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["game"]

	err := s.DeleteGame(gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := json.NewEncoder(w)
	out.Encode(true)
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

// Leave an existing game
func (s *Server) leaveGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["game"]

	playerID := r.FormValue("player")

	game, err := s.GetGame(gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if game.HomeID == playerID {
		game.HomeID = ""
		game.HomeReady = false
	} else if game.AwayID == playerID {
		game.AwayID = ""
		game.AwayReady = false
	}

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
	gameID := vars["game"]

	playerID := r.FormValue("player")

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
