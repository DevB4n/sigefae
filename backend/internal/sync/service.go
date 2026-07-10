package sync

import (
	"encoding/json"
	"os"
	"time"
)

const stateFile = "internal/sync/state.json"

type State struct {
	LastSync time.Time `json:"last_sync"`
}

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) LastSync() time.Time {
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return time.Time{}
	}

	var state State

	if err := json.Unmarshal(data, &state); err != nil {
		return time.Time{}
	}

	return state.LastSync
}

func (s *Service) Update(last time.Time) error {
	state := State{
		LastSync: last,
	}

	data, err := json.MarshalIndent(state, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(stateFile, data, 0644)
}