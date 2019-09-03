package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp/da-dance-api/models"
	// _ "github.com/lib/pq"
)

// AllocationRequest holds allocation data.
type AllocationRequest struct {
	Player string `json:"player"`
	Count  int    `json:"count"`
}

// AssignAllocations reads allocations from the database, compares it to a
// running list, and assigns a player to them.
func (s *Server) AssignAllocations(count int) ([]string, error) {
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
		exists := false
		for _, existing := range existingAllocations {
			if id == existing {
				exists = true
				break
			}
		}

		if !exists && len(assigned) < count {
			_, err := s.CreateAllocation(models.Allocation{ID: id})
			if err != nil {
				return nil, err
			}
			assigned = append(assigned, id)
		}
	}

	s.logger.Info("assigned allocations", "num_allocations", len(assigned))
	return assigned, nil
}

// Get allocations
func (s *Server) getAllocationsHandler(w http.ResponseWriter, r *http.Request) {
	c := r.FormValue("count")
	count, err := strconv.Atoi(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allocations, err := s.AssignAllocations(count)
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

	err := s.DeleteAllocation(models.Allocation{ID: allocationID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.logger.Info("killed allocation", "id", allocationID)

	out := json.NewEncoder(w)
	out.Encode(allocationID)
}

// Clear an allocation
func (s *Server) clearAllocationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	allocationID := vars["allocation"]

	err := s.DeleteAllocation(models.Allocation{ID: allocationID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.logger.Info("released claim", "id", allocationID)

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
