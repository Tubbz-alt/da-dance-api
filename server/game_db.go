package server

import (
	"database/sql"

	"github.com/hashicorp/da-dance-api/models"
)

// GetGames gets all games
func (s *Server) GetGames() ([]models.Game, error) {
	games := []models.Game{}
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

// DeleteGame deletes an existing game
func (s *Server) DeleteGame(id string) error {
	// Delete game
	deleteQuery, err := s.database.Prepare(
		`DELETE FROM games
		WHERE id = $1`)
	if err != nil {
		return err
	}

	_, err = deleteQuery.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
