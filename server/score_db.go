package server

import (
	"database/sql"

	"github.com/hashicorp/da-dance-api/models"
)

// GetScore gets a score by player ID
func (s *Server) GetScore(id string) (models.Score, error) {
	score := models.Score{
		Player: id,
	}

	scoreQuery, err := s.database.PrepareNamed(
		`SELECT * FROM scores
		WHERE player = :player`)
	if err != nil {
		s.logger.Error("Prepare score query", "error", err)
	}

	err = scoreQuery.Get(&score, score)
	if err == sql.ErrNoRows {
		return score, err
	}

	if err != nil {
		s.logger.Error("Query score", "error", err)
		return score, err
	}

	return score, nil
}

func (s *Server) updateScore(score models.Score) error {
	query, err := s.database.PrepareNamed(
		`UPDATE scores
		SET
			game = :game, 
			points = :points
		WHERE player = :player
		RETURNING *`)
	err = query.Get(&score, score)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) insertScore(score models.Score) error {
	query, err := s.database.PrepareNamed(
		`INSERT INTO scores (player, game, points)
		VALUES(:player, :game, :points)
		RETURNING *`)
	err = query.Get(&score, score)
	if err != nil {
		return err
	}
	return nil
}

// CreateScore creates a scores
func (s *Server) CreateScore(score models.Score) (models.Score, error) {
	currentScore, err := s.GetScore(score.Player)
	if err == sql.ErrNoRows {
		err = s.insertScore(score)
	} else if err == nil && currentScore.Points < score.Points {
		err = s.updateScore(score)
	}

	if err != nil {
		return score, err
	}

	return score, nil
}

// GetAllScores creates a scores
func (s *Server) GetAllScores() ([]models.Score, error) {
	scores := []models.Score{}
	err := s.database.Select(&scores, "SELECT * FROM scores ORDER BY points DESC LIMIT 10")
	if err == sql.ErrNoRows {
		s.logger.Error("No scores found")
		return scores, nil
	}

	if err != nil {
		s.logger.Error("Get scores", "error", err)
		return scores, err
	}

	return scores, nil
}

// GetScores creates a scores
func (s *Server) GetScores(game string) ([]models.Score, error) {
	params := map[string]interface{}{
		"game": game,
	}
	scores := []models.Score{}
	query, err := s.database.PrepareNamed(`SELECT * FROM scores WHERE game = :game`)
	err = query.Select(&scores, params)
	if err == sql.ErrNoRows {
		s.logger.Error("No scores found")
		return scores, nil
	}

	if err != nil {
		s.logger.Error("Get scores", "error", err)
		return scores, err
	}

	return scores, nil
}
