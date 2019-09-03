package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/da-dance-api/models"
)

// ScoreRequest holds score data.
type ScoreRequest struct {
	Player string `json:"player"`
	Points int    `json:"points"`
}

func (s *Server) getAllScoresHandler(w http.ResponseWriter, r *http.Request) {
	scores, err := s.GetAllScores()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := json.NewEncoder(w)
	out.Encode(scores)
}

func (s *Server) getScoresHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["game"]
	scores, err := s.GetScores(gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := json.NewEncoder(w)
	out.Encode(scores)
}

func (s *Server) createScoreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["game"]

	var sr ScoreRequest
	err := json.NewDecoder(r.Body).Decode(&sr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	score := models.Score{
		Game:   gameID,
		Player: sr.Player,
		Points: sr.Points,
	}

	insertedScore, err := s.CreateScore(score)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out := json.NewEncoder(w)
	out.Encode(insertedScore)
}
