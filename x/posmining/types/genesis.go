package types

import (
	"time"
)

// GenesisState - all posmining state that must be provided at genesis
type GenesisState struct {
	Records []Posmining
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(records []Posmining) GenesisState {
	return GenesisState{
		Records: records,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	startDate := time.Date(2019, 9, 1, 0, 0, 0, 0, time.UTC)

	return GenesisState{
		Records: make([]Posmining, 0),
	}
}

// ValidateGenesis validates the posmining genesis parameters
func ValidateGenesis(data GenesisState) error {
	return nil
}
