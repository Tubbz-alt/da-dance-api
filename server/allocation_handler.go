package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hashicorp/da-dance-api/models"
	"github.com/gorilla/mux"
	// _ "github.com/lib/pq"
)

// AssignAllocations reads allocations from the database, compares it to a
// running list, and assigns a player to them.
func (s *Server) AssignAllocations(playerID string, count int) ([]string, error) {
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
		if !exists && len(assigned) < count {
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
	count := r.FormValue("count")

	allocationCount, err := strconv.Atoi(count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allocations, err := s.AssignAllocations(player, allocationCount)
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

	if err := s.nomad.StopAllocation(allocationID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allocation, err := s.DeleteAllocation(models.Allocation{ID: allocationID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.logger.Info("released claim", "id", allocationID, "player", allocation.Player)

	out := json.NewEncoder(w)
	out.Encode(allocationID)
}

// Clear an allocation
func (s *Server) clearAllocationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	allocationID := vars["allocation"]

	allocation, err := s.DeleteAllocation(models.Allocation{ID: allocationID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.logger.Info("released claim", "id", allocationID, "player", allocation.Player)

	out := json.NewEncoder(w)
	out.Encode(allocationID)
}

// Clear all allocation
func (s *Server) clearAllocationsHandler(w http.ResponseWriter, r *http.Request) {
	err := s.DeleteAllocations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.logger.Info("released all claims")

	out := json.NewEncoder(w)
	out.Encode(true)
}
