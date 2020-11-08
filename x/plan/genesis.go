package plan

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
}

// New Genesis
func NewGenesisState() GenesisState {
	return GenesisState{}
}

// Validation
func ValidateGenesis(data GenesisState) error {
	return nil
}

// Default state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// Init
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// Export
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{}
}
