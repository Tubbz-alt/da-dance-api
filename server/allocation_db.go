package server

import (
	"database/sql"
	"errors"

	"github.com/hashicorp/da-dance-api/models"
)

const (
	defaultAllocationBatchSize = 10
)

// GetAllocations gets all allocations
func (s *Server) GetAllocations() (map[string]string, error) {
	allocation := map[string]string{}
	var tempAllocs []models.Allocation
	err := s.database.Select(&tempAllocs, "SELECT * FROM allocations")
	if err == sql.ErrNoRows {
		s.logger.Error("No allocations found")
		return allocation, nil
	}

	if err != nil {
		s.logger.Error("Get allocations", "error", err)
		return allocation, err
	}

	for _, alloc := range tempAllocs {
		allocation[alloc.ID] = alloc.Player
	}

	return allocation, nil
}

// CreateAllocation creates an allocation
func (s *Server) CreateAllocation(allocation models.Allocation) (models.Allocation, error) {
	query, err := s.database.PrepareNamed(
		`INSERT INTO allocations (id, player)
		VALUES(:id, :player)
		RETURNING *`)
	if err != nil {
		return allocation, err
	}

	err = query.Get(&allocation, allocation)
	if err != nil {
		return allocation, err
	}

	return allocation, nil
}

// DeleteAllocation removes an allocation to demonstrate it is released
func (s *Server) DeleteAllocation(allocation models.Allocation) (models.Allocation, error) {
	query, err := s.database.PrepareNamed(
		`DELETE FROM allocations WHERE id=:id RETURNING *`)

	if err != nil {
		return allocation, err
	}

	err = query.Get(&allocation, allocation)
	if err == sql.ErrNoRows {
		return allocation, errors.New("Allocation not found")
	}

	if err != nil {
		return allocation, err
	}

	return allocation, nil
}

// DeleteAllocations removes all allocations
func (s *Server) DeleteAllocations() error {
	query, err := s.database.Prepare(
		`TRUNCATE allocations`)
	if err != nil {
		return err
	}

	_, err = query.Exec()
	if err != nil {
		return err
	}

	return nil
}
