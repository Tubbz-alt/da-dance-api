package nomad

import (
	"github.com/hashicorp/nomad/api"
)

// Nomad has a client and allocation endpoint
type Nomad struct {
	Client      *api.Client
	Allocations *api.Allocations
	Assignments map[string]string
}

// Connect creates a nomad client and initializes with a list of running allocations
func Connect() (*Nomad, error) {
	client := Nomad{}
	if err := client.DefaultClient(); err != nil {
		return nil, err
	}
	client.Assignments = map[string]string{}
	return &client, nil
}

// DefaultClient creates a new nomad client with default configuration
func (n *Nomad) DefaultClient() error {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	n.Client = client
	n.Allocations = client.Allocations()
	return nil
}

// StopAllocation stops an allocation based on ID
func (n *Nomad) StopAllocation(id string) error {
	alloc := &api.Allocation{ID: id}
	queryOptions := &api.QueryOptions{
		AllowStale: true,
	}
	n.Allocations.Stop(alloc, queryOptions)
	return nil
}

// GetAssignableAllocations returns a list of currently running allocations
func (n *Nomad) GetAssignableAllocations() ([]string, error) {
	alloc := []string{}
	queryOptions := &api.QueryOptions{
		AllowStale: true,
	}
	allocationList, _, err := n.Allocations.List(queryOptions)
	if err != nil {
		return nil, err
	}
	runningAllocations := searchForRunning(allocationList)
	for _, allocation := range runningAllocations {
		if _, ok := n.Assignments[allocation.ID]; !ok {
			alloc = append(alloc, allocation.ID)
		}
	}
	return alloc, nil
}

func searchForRunning(allocationList []*api.AllocationListStub) []*api.AllocationListStub {
	runningAllocations := []*api.AllocationListStub{}
	for _, allocation := range allocationList {
		if allocation.ClientStatus == "running" {
			runningAllocations = append(runningAllocations, allocation)
		}
	}
	return runningAllocations
}
