package server

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/eveld/ddr-api/models"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	defaultAllocationBatchSize = 10
)

// GetAllocations gets all allocations
func (s *Server) GetAllocations() (map[string]bool, error) {
	allocation := map[string]bool{}
	var tempAllocs []string
	err := s.database.Select(&tempAllocs, "SELECT id FROM allocations")
	if err == sql.ErrNoRows {
		s.logger.Error("No allocations found")
		return allocation, nil
	}

	if err != nil {
		s.logger.Error("Get allocations", "error", err)
		return allocation, err
	}

	for _, alloc := range tempAllocs {
		allocation[alloc] = true
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

func (s *Server) assignAllocations(playerID string) ([]string, error) {
	assigned := []string{}
	existingAllocations, err := s.GetAllocations()
	if err != nil {
		return nil, err
	}
	runningAllocations, err := s.nomad.GetRunningAllocations()
	if err != nil {
		return nil, err
	}

	for _, id := range runningAllocations {
		_, exists := existingAllocations[id]
		if !exists && len(assigned) < defaultAllocationBatchSize {
			_, err := s.CreateAllocation(models.Allocation{ID: id, Player: playerID})
			if err != nil {
				return nil, err
			}
			assigned = append(assigned, id)
		}
	}
	
	s.logger.Info("assigned allocations", "num_allocations", len(assigned), "player", playerID)
	return assigned, nil
}

// Get allocations
func (s *Server) getAllocationsHandler(w http.ResponseWriter, r *http.Request) {
	player := r.FormValue("player")
	allocations, err := s.assignAllocations(player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out := json.NewEncoder(w)
	out.Encode(allocations)
}

// Stop an allocation
func (s *Server) stopAllocationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	allocationID := vars["allocation"]

	out := json.NewEncoder(w)
	out.Encode(allocationID)
}
