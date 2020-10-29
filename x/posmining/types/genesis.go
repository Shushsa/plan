package types

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
	return GenesisState{
		Records: make([]Posmining, 0),
	}
}

// ValidateGenesis validates the posmining genesis parameters
func ValidateGenesis(data GenesisState) error {
	return nil
}
