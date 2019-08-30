package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/eveld/ddr-api/models"
)

func (s *Server) getScoresHandler(w http.ResponseWriter, r *http.Request) {
	scores, err := s.GetScores()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := json.NewEncoder(w)
	out.Encode(scores)
}

func (s *Server) createScoreHandler(w http.ResponseWriter, r *http.Request) {
	player := r.FormValue("player")
	game := r.FormValue("game")
	p := r.FormValue("points")

	points, err := strconv.Atoi(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	score := models.Score{
		Player: player,
		Game:   game,
		Points: points,
	}

	insertedScore, err := s.CreateScore(score)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out := json.NewEncoder(w)
	out.Encode(insertedScore)
}
