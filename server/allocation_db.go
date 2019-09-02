package server

import (
	"database/sql"

	"github.com/hashicorp/da-dance-api/models"
)

const (
	defaultAllocationBatchSize = 10
)

// GetAllocations gets all allocations
func (s *Server) GetAllocations() ([]string, error) {
	var allocations []string
	err := s.database.Select(&allocations, "SELECT id FROM allocations")
	if err == sql.ErrNoRows {
		s.logger.Error("No allocations found")
		return allocations, nil
	}

	if err != nil {
		s.logger.Error("Get allocations", "error", err)
		return allocations, err
	}

	return allocations, nil
}

// CreateAllocation creates an allocation
func (s *Server) CreateAllocation(allocation models.Allocation) (models.Allocation, error) {
	query, err := s.database.PrepareNamed(
		`INSERT INTO allocations (id)
		VALUES(:id)
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
func (s *Server) DeleteAllocation(allocation models.Allocation) error {
	query, err := s.database.PrepareNamed(
		`DELETE FROM allocations WHERE id=:id RETURNING *`)
	if err != nil {
		return err
	}

	_, err = query.Exec(allocation)
	if err != nil {
		return err
	}

	return nil
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
