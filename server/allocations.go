package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	defaultAllocationBatchSize = 3
)

func (s *Server) AssignAllocations(playerID string) ([]string, error) {
	var assigned []string

	allocations, err := s.nomad.GetAssignableAllocations()
	if err != nil {
		return nil, err
	}

	if len(allocations) > defaultAllocationBatchSize {
		assigned = allocations[:defaultAllocationBatchSize]
	} else {
		assigned = allocations
	}

	for _, id := range assigned {
		s.nomad.Assignments[id] = playerID
	}
	s.logger.Info("assigned allocations", "num_allocations", len(assigned), "player", playerID)
	return assigned, nil
}

// Get allocations
func (s *Server) getAllocationsHandler(w http.ResponseWriter, r *http.Request) {
	player := r.FormValue("player")
	allocations, err := s.AssignAllocations(player)
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
